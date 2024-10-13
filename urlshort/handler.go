package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		url, exists := pathsToUrls[req.URL.Path]

		if !exists {
			fallback.ServeHTTP(res, req)
			return
		}

		http.Redirect(res, req, url, http.StatusSeeOther)
	}
}

type pathToUrl struct {
	Path string
	Url  string
}

func GetHandler(filename *string, data []byte, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	extension := filepath.Ext(*filename)
	var handler http.HandlerFunc
	var err error
	switch extension {
	case ".yaml", ".yml":
		handler, err = YAMLHandler(data, fallback)
	case ".json":
		handler, err = JSONHandler(data, fallback)
	default:
		err = fmt.Errorf("unknown file extension '%s'", extension)
	}
	if err != nil {
		return nil, err
	}
	return handler, err
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(paths)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYAML(yml []byte) ([]pathToUrl, error) {
	var res []pathToUrl
	err := yaml.Unmarshal(yml, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseJSON(data)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(paths)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseJSON(data []byte) ([]pathToUrl, error) {
	var res []pathToUrl
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func buildMap(paths []pathToUrl) map[string]string {
	res := make(map[string]string, len(paths))
	for _, mapping := range paths {
		res[mapping.Path] = mapping.Url
	}
	return res
}
