package urlshort

import (
	"fmt"
	"net/http"
)

// PathMap represents name-url mappings
type PathMap = map[string]string

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

		http.Redirect(w, r, redirectURL, 301)
	}
}
