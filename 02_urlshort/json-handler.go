package urlshort

import (
	"encoding/json"
	"net/http"
)

func buildJSONMap(JSONMap []pathURL) PathMap {
	urlMap := make(PathMap)
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
