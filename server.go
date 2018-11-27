package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	port = ":8080"
)

func runServer() {
	fs := http.FileServer(http.Dir(parrotDirectory))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", parrotViewHandler)
	fmt.Printf("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(port, nil))

}

func parrotViewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Parrot list</h1>")
	files, err := ioutil.ReadDir(parrotDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Fprintf(w, "<img src='/static/%v'></img><h4>%v</h4>", f.Name(), f.Name())
	}
}
