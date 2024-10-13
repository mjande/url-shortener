package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mjande/url-shortener/urlshort"
)

func main() {
	filename := flag.String("file", "default.yaml", "A YAML file in format\n- path: '/github'\n  url: 'http://github.com'\n")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Get YAML from file
	yaml, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Printf("ERROR: Unable to parse file %s\n", *filename)
		os.Exit(1)
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
