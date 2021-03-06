package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"some-repo-url/urlshort"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	// Build the MapHandler using the mux as the fallback
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	customYamlPointer := flag.String("yaml", "", "path to yaml file")
	customJSONPointer := flag.String("json", "", "path to json file")
	flag.Parse()

	jsonHandler := setupHandlers(customYamlPointer, customJSONPointer, mapHandler)

	finalHandler := urlshort.PgHandler(jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", finalHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func setupHandlers(
	customYamlPointer *string,
	customJSONPointer *string,
	 mapHandler http.HandlerFunc,
) http.HandlerFunc {
	var yaml string
	var json string

	if (*customYamlPointer != "") {
		fmt.Println("using yaml from file path", *customYamlPointer)

		data, err := ioutil.ReadFile(*customYamlPointer)
		check(err)

		yaml = string(data)
	} else {
		fmt.Println("using default yaml")

		yaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	}

	if (*customJSONPointer != "") {
		fmt.Println("using JSON from file path", *customJSONPointer)

		data, err := ioutil.ReadFile(*customJSONPointer)
		check(err)

		json = string(data)
	} else {
		fmt.Println("using default json")

		json = `[{"path":"/json","url":"https://github.com/gophercises/urlshort"}]`
	}

	// Build the YAMLHandler using the mapHandler
	// as the fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	check(err)

	// Build the JSONHandler using the YAMLHandler
	// as the fallback
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	check(err)

	return jsonHandler
}
