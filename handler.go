package main

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		url, ok := pathsToUrls[path]
		if ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
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

	sliceOfPathAndURL, err := createSliceOfPathAndURL(yml)
	if err != nil {
		return nil, err
	}

	mapOfPathAndURL := createMapOfPathAndURL(sliceOfPathAndURL)

	return MapHandler(mapOfPathAndURL, fallback), nil
}

func createSliceOfPathAndURL(yml []byte) ([]pathAndURL, error) {
	var slicePU []pathAndURL
	err := yaml.Unmarshal(yml, &slicePU)
	return slicePU, err
}

func createMapOfPathAndURL(slicePU []pathAndURL) map[string]string {
	mapPU := make(map[string]string)
	for _, a := range slicePU {
		mapPU[a.Path] = a.URL
	}
	return mapPU
}

type pathAndURL struct {
	Path string
	URL  string
}
