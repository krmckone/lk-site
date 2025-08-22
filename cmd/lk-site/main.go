package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/krmckone/lk-site/internal/templating"
	"github.com/krmckone/lk-site/internal/utils"
)

func main() {
	assetsPath := flag.String("assets-path", "assets", "Path to the assets directory")
	configsPath := flag.String("configs-path", "configs", "Path to the configs directory")
	buildPath := flag.String("build-path", "configs", "Path to the build directory")
	flag.Parse()
	runtime := utils.RuntimeConfig{
		AssetsPath:  *assetsPath,
		ConfigsPath: *configsPath,
		BuildPath:   *buildPath,
	}
	if err := templating.TemplateSite(runtime); err != nil {
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
