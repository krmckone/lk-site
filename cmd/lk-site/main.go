package main

import (
	"bytes"

	"github.com/krmckone/ksite/internal/config"
	"github.com/krmckone/ksite/internal/templater"
	"github.com/krmckone/ksite/internal/utils"
	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// For each md file:
// 1. Read in
// 2. Template any custom variables
// 3. Render to HTML
func main() {
	utils.Clean("build")
	config := config.ReadConfig("configs/config.yml")

	buf := new(bytes.Buffer)
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

	// Do main page content
	// Read in
	md := utils.ReadFile("assets/index.md")
	// Template
	md, err := templater.Run(md, config.Template.Params)
	if err != nil {
		panic(err)
	}
	// Render to HTML
	if err := gm.Convert(md, buf); err != nil {
		panic(err)
	}
	config.Template.Params["main_content"] = buf.String()
	buf.Reset()

	// Do header and topnav content
	// Read in
	md = utils.ReadFile("assets/topnav.md")
	// Template
	md, err = templater.RunNav(md, config.Template.Params)
	if err != nil {
		panic(err)
	}
	// Render to HTML
	if err := gm.Convert(md, buf); err != nil {
		panic(err)
	}
	config.Template.Params["topnav"] = buf.String()
	buf.Reset()
	// Read in
	md = utils.ReadFile("assets/header.md")
	// Template
	md, err = templater.Run(md, config.Template.Params)
	if err != nil {
		panic(err)
	}
	// Render to HTML
	if err := gm.Convert(md, buf); err != nil {
		panic(err)
	}
	config.Template.Params["header"] = buf.String()
	buf.Reset()

	// Do footer content
	// Read in
	md = utils.ReadFile("assets/footer.md")
	// Template
	md, err = templater.Run(md, config.Template.Params)
	if err != nil {
		panic(err)
	}
	// Render to HTML
	if err := gm.Convert(md, buf); err != nil {
		panic(err)
	}
	config.Template.Params["footer"] = buf.String()
	buf.Reset()

	// Build from the base page and rendered MD content
	// Read in
	basePage := utils.ReadFile("assets/base_page.html")
	// Template
	md, err = templater.Run(basePage, config.Template.Params)
	if err != nil {
		panic(err)
	}
	utils.Mkdir("build")
	utils.WriteFile("build/index.html", md)
}
