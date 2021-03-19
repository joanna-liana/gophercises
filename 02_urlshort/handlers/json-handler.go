package urlshort_handlers

import (
	"encoding/json"
	"net/http"
)

// JSONHandler is analagous to YAMLHandler
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedJSON []PathURL

	err := json.Unmarshal(jsonBytes, &parsedJSON)

	if err != nil {
		return nil, err
	}

	pathMap := BuildPathMapFromPathURLs(parsedJSON)

	return MapHandler(pathMap, fallback), err
}
