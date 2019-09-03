package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"

	urlshort "github.com/WeirdMagician/urlshortner"
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

	fmt.Print("started")

	db, err := bolt.Open("C:\\Users\\kh2398\\go\\src\\github.com\\WeirdMagician\\urlshortner\\main\\my.db", 0666, nil)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("myBucket"))
		if err != nil {
			log.Fatal(err)
		}
		err = b.Put([]byte("/bolt"), []byte("https://github.com/boltdb/bolt"))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	bolthandler := urlshort.BoltHandler(mux)

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, bolthandler)

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
