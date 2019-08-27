package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/WeirdMagician/urlshort"
)

func main() {
	defaultyml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution`

	defaultjson := `[{"path": "/kony","url": "https://kony.com"},{"path": "/indeed","url": "https://indeed.com"}]`
	filenameYML := flag.String("yaml-file", defaultyml, "Give the yaml file has input with below format")
	filenameJSON := flag.String("json-file", defaultjson, "Give the yaml file has input with below format")

	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamld := check(filenameYML, defaultyml)
	jsond := check(filenameJSON, defaultjson)

	yamlHandler, err := urlshort.YAMLHandler([]byte(yamld), []byte(jsond), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", yamlHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func check(fn *string, de string) string {
	if *fn == de {
		return de
	}
	b, err := ioutil.ReadFile(*fn)
	if err != nil {
		log.Fatalln("Error :", err)
	}
	return string(b)

}
