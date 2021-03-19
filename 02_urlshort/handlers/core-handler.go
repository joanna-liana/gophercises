package urlshort_handlers

import (
	"fmt"
	"net/http"
)

// PathMap represents name-url mappings
type PathMap = map[string]string

// PathURL represents the struct used for (un)marshalling data
type PathURL struct {
	Path string
	Url string
}


// BuildPathMapFromPathURLs does what the name says
func BuildPathMapFromPathURLs(pathURLs []PathURL) PathMap {
	urlMap := make(PathMap)
	for _, pu := range pathURLs {
		urlMap[pu.Path] = pu.Url
	}
	return urlMap
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls PathMap, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := pathsToUrls[r.URL.String()]
		fmt.Println(redirectURL)

		if redirectURL == "" {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
