package main

import (
	"github.com/krmckone/ksite/internal/config"
	"github.com/krmckone/ksite/internal/templater"
	"github.com/krmckone/ksite/internal/utils"
	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	utils.Clean("build")
	config := config.ReadConfig("configs/config.yml")

	gm := goldmark.New(
		attributes.Enable,
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)

	if err := templater.BuildSite(gm, &config); err != nil {
		panic(err)
	}
}
