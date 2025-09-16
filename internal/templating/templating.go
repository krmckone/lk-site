package templating

import (
	"bytes"
	"fmt"
	gohtml "html"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/krmckone/lk-site/internal/config"
	"github.com/krmckone/lk-site/internal/page"
	"github.com/krmckone/lk-site/internal/utils"
	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// BuildSite is for building the site. This includes templating HTML with markdown and
// putting images in the expected locations in the output
func TemplateSite(runtime utils.RuntimeConfig) error {
	utils.SetupBuild(runtime)

	c, err := config.ReadConfig(runtime)
	if err != nil {
		return err
	}

	pages, err := getAssetPages(runtime, "", c.Template.Params)
	if err != nil {
		return err
	}

	assetTemplatePaths := utils.GetBasePageFiles(runtime)

	// The main content of each page can refer to other templates that are defined separately,
	// so we need to template the main content as well against any component templates. We'll
	// pass this list of component files to the setupPageParams function so that it can
	// template the main content against them
	componentFiles, err := utils.GetComponentFiles(runtime)
	if err != nil {
		return err
	}

	tmpl := template.New("base_page.html")
	tmpl, err = tmpl.Funcs(runtime.TemplateFuncs).ParseFiles(assetTemplatePaths...)
	if err != nil {
		log.Printf("Error parsing files: %s, %s", assetTemplatePaths, err)
		return err
	}

	gm := newGoldmark()
	for _, page := range pages {
		mdBuffer := bytes.Buffer{}
		if err := gm.Convert(page.Content, &mdBuffer); err != nil {
			return err
		}

		pageParams, err := setupPageParams(
			runtime,
			componentFiles,
			c,
			mdBuffer.String(),
		)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(
			filepath.Dir(page.BuildPath),
			os.ModePerm,
		); err != nil {
			return err
		}
		file, err := os.Create(page.BuildPath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.ExecuteTemplate(file, "base_page.html", pageParams); err != nil {
			log.Printf("Error executing template: %s, %s", pageParams, err)
			return err
		}
	}

	return nil
}

func setupPageParams(runtime utils.RuntimeConfig, componentFiles []string, config config.Config, mainContent string) (map[string]interface{}, error) {
	pageParams := map[string]interface{}{}
	for k, v := range config.Template.Params {
		pageParams[k] = template.HTML(v.(string))
	}
	for k, v := range config.Env.Params {
		pageParams[k] = v.(string)
	}
	mainContentTemplate, err := template.Must(
		template.New("main_content").Funcs(runtime.TemplateFuncs).Parse(gohtml.UnescapeString(mainContent)),
	).ParseFiles(
		componentFiles...,
	)
	if err != nil {
		return nil, err
	}
	mainContentBuffer := bytes.Buffer{}
	if err := mainContentTemplate.ExecuteTemplate(&mainContentBuffer, "main_content", pageParams); err != nil {
		log.Printf("Error executing template: %s, %s", pageParams, err)
		return nil, err
	}
	pageParams["main_content"] = template.HTML(mainContentBuffer.String())
	pageParams["title"] = config.Template.Params["title"].(string)
	return pageParams, nil
}

func newGoldmark() goldmark.Markdown {
	return goldmark.New(
		attributes.Enable,
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(), // Lets you use {.att } syntax to add attributes to HTML output
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(), // For embedding our own HTML in markdown files. Without this, goldmark will hide the HTML content in the output
		),
	)
}

// a recursive function that returns all the pages in the directory
// and all subdirectories. The path is a parameter because this
// function is called recursively to get all the pages in the site underneath
// "assets/pages". The first call of this function should be with the empty string
// as path which represents the root of assets/pages
func getAssetPages(runtime utils.RuntimeConfig, path string, params map[string]interface{}) ([]page.Page, error) {
	baseAssetPath := utils.MakePath(filepath.Join(runtime.AssetsPath, "pages"))
	fullAssetPath := path
	if !strings.HasPrefix(path, baseAssetPath) {
		fullAssetPath = filepath.Join(baseAssetPath, path)
	}
	pages := []page.Page{}

	files, err := os.ReadDir(fullAssetPath)
	if err != nil {
		return pages, err
	}

	for _, file := range files {
		if file.IsDir() {
			subPages, err := getAssetPages(runtime, filepath.Join(fullAssetPath, file.Name()), params)
			if err != nil {
				return pages, err
			}
			pages = append(pages, subPages...)
		}
		// We only allow markdown files to be pages in this context
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		// Read markdown content
		content, err := os.ReadFile(filepath.Join(fullAssetPath, file.Name()))
		if err != nil {
			return pages, err
		}

		// Mark where the output for this page should be written
		title := strings.TrimSuffix(file.Name(), ".md")
		buildPath := filepath.Join(
			strings.ReplaceAll(fullAssetPath, baseAssetPath, utils.MakePath("build")),
			fmt.Sprintf("%s.html", title),
		)
		page := page.Page{
			Title:     title,
			Content:   content,
			Params:    params,
			AssetPath: fullAssetPath,
			BuildPath: buildPath,
		}
		pages = append(pages, page)
	}

	return pages, nil
}
