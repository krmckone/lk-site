package main

import (
	"fmt"
	"log"
	"net/http"
)

// Use to test the output of the static site generator at localhost:3000
func main() {
	port := "8080"
	directory := "./build"
	log.Printf("Serving %s on HTTP port: %s\n", directory, port)
	http.Handle("/", http.FileServer(http.Dir(directory)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
