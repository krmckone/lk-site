package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/krmckone/ksite/internal/templater"
	"github.com/krmckone/ksite/internal/utils"
)

func main() {
	resp, err := utils.GetSteamOwnedGames()
	if err != nil {
		panic(err)
	}
	fmt.Println("TESTING")
	fmt.Println(utils.GetTopFiftySteamDeckGames(resp))
	fmt.Println("TESTING")
	if err := templater.BuildSite(); err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "server" {
		port := args[1]
		directory := "./build"
		log.Printf("Serving %s on HTTP port: %s\n", directory, port)
		http.Handle("/", http.FileServer(http.Dir(directory)))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}

}
