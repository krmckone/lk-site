package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/krmckone/lk-site/internal/templating"
)

func main() {
	if err := templating.TemplateSite(); err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "server" {
		port := args[1]
		serveDir := "./build"
		log.Printf("Serving %s on HTTP port: %s\n", serveDir, port)
		http.Handle("/", http.FileServer(http.Dir(serveDir)))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}

}
