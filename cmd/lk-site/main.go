package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/krmckone/ksite/internal/config"
	"github.com/krmckone/ksite/internal/preprocessor"
	"github.com/krmckone/ksite/internal/utils"
)

func main() {
	utils.Clean("build")
	config := config.ReadConfig("configs/config.yml")
	md := utils.ReadFile("assets/index.md")
	md = preprocessor.Run(md, config.Template.Params)
	htmlOpts := html.RendererOptions{
		CSS:   config.Template.Styles.SheetURL,
		Flags: html.CommonFlags | html.CompletePage,
	}
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.Attributes)
	renderer := html.NewRenderer(htmlOpts)
	output := markdown.ToHTML(md, parser, renderer)
	utils.Mkdir("build")
	utils.WriteFile("build/index.html", output)
}
