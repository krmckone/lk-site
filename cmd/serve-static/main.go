package main

import (
	"net/http"
)

// Use to test the output of the static site generator at localhost:3000
func main() {
	http.Handle("/", http.FileServer(http.Dir("./build")))
	http.ListenAndServe(":3000", nil)
}
