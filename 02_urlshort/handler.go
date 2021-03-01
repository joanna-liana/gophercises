package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathMap = map[string]string

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls pathMap, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := pathsToUrls[r.URL.String()]
		fmt.Println(redirectURL)

		if redirectURL == "" {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, redirectURL, 301)
	}
}

func buildYamlMap(yamlMap []map[string]string) pathMap {
	urlMap := make(pathMap)
	for _, dict := range yamlMap {
		urlMap[dict["path"]] = dict["url"]
	}
	return urlMap
}


// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var parsedYaml []pathMap

	err := yaml.Unmarshal(yml, &parsedYaml)

	if err != nil {
		return nil, err
	}

	pathMap := buildYamlMap(parsedYaml)

	return MapHandler(pathMap, fallback), err
}

func buildJSONMap(JSONMap []pathURL) pathMap {
	urlMap := make(pathMap)
	for _, entry := range JSONMap {
		urlMap[entry.Path] = entry.Url
	}
	return urlMap
}

type pathURL struct {
	Path string
	Url string
}

// JSONHandler is analagous to YAMLHandler
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedJSON []pathURL

	err := json.Unmarshal(jsonBytes, &parsedJSON)

	if err != nil {
		return nil, err
	}

	pathMap := buildJSONMap(parsedJSON)


	return MapHandler(pathMap, fallback), err
}

// TODO:
// Build a Handler that doesn’t read from a map but instead reads from a database. Whether you use BoltDB, SQL, or something else is entirely up to you.
