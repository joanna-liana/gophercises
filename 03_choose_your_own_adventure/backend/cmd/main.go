package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"cyoa/story"
)

// TODO: CLI port
func startStory(parsedStory story.ParsedStory) http.Handler {
	intro := story.GetStoryIntro(parsedStory)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(intro)
	})
}

func chooseNextArc(parsedStory story.ParsedStory) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		arcName := strings.TrimPrefix(req.URL.Path, "/next/")

		err, arc := story.GetArcByName(parsedStory, arcName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(arc)
	})
}

func main() {
	parsedStory := story.GetParsedStory("story/gopher.json")

	http.Handle("/start", startStory(parsedStory))
	http.Handle("/next/", chooseNextArc(parsedStory))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
