package urlshorter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/boltdb/bolt"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	return func(w http.ResponseWriter, r *http.Request) {
		redirecturl, ok := pathsToUrls[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, redirecturl, http.StatusFound)
		}
	}
}

// BoltHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the db) to their corresponding URL (values
// that each key in the db points to, in []byte format).
// If the path is not provided in the db, then the fallback
// http.Handler will be called instead.
func BoltHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("hey\n")

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recoverd", r)
			}
		}()

		db, err := bolt.Open("C:\\Users\\kh2398\\go\\src\\github.com\\WeirdMagician\\urlshortner\\main\\my.db", 0666, nil)
		fmt.Print("redirecturl", "redirecturl\n")

		if err != nil {
			fmt.Print(err)
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recoverd", r)
			}
			db.Close()
		}()

		var redirecturl string
		err = db.View(func(tx *bolt.Tx) error {
			fmt.Fprint(os.Stdout, tx.Bucket([]byte("myBucket")).Get([]byte(r.URL.Path)), "tredirecturl\n")

			redirecturl = string(tx.Bucket([]byte("myBucket")).Get([]byte(r.URL.Path)))
			fmt.Fprint(os.Stdout, redirecturl, "redirecturl\n")
			return nil
		})
		if err != nil {
			fmt.Print(err)
		}
		if redirecturl == "" {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, string(redirecturl), http.StatusFound)
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
func YAMLHandler(yml []byte, jsond []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...

	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	parsedJSON, err := parseJSON(jsond)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(append(parsedYaml, parsedJSON...))

	return MapHandler(pathMap, fallback), nil
}

type link struct {
	Path string
	URL  string
}

func parseYAML(yml []byte) ([]link, error) {
	var list []link
	err := yaml.Unmarshal(yml, &list)
	return list, err
}

func parseJSON(jsond []byte) ([]link, error) {
	var list []link
	err := json.Unmarshal(jsond, &list)
	return list, err
}

func buildMap(parsedYaml []link) map[string]string {
	s := make(map[string]string)
	for _, list := range parsedYaml {
		s[list.Path] = list.URL
	}
	return s
}
