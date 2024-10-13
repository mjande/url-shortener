package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mjande/url-shortener/urlshort"
)

func main() {
	filename := flag.String("file", "default.yaml", "A JSON or YAML file")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Get data from file
	data, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Printf("Unable to parse file %s\n", *filename)
		os.Exit(1)
	}

	// Build the get handler using the mapHandler as the
	// fallback
	handler, err := urlshort.GetHandler(filename, data, mapHandler)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
