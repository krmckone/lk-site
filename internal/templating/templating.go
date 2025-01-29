package templating

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/krmckone/lk-site/internal/config"
	"github.com/krmckone/lk-site/internal/page"
	"github.com/krmckone/lk-site/internal/utils"
	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// BuildSite is for building the site. This includes templating HTML with markdown and
// putting images in the expected locations in the output
func TemplateSite() error {
	utils.SetupBuild()

	c, err := config.ReadConfig("configs/config.yml")
	if err != nil {
		return err
	}

	pages, err := getAssetPages("", c.Template.Params)
	if err != nil {
		return err
	}

	funcs := getTemplateFuncs()
	tmpl := template.New(filepath.Join(utils.MakePath("assets"), "base_page.html"))
	tmpl, err = tmpl.Funcs(funcs).ParseFiles(
		filepath.Join(utils.MakePath("assets"), "base_page.html"),
		filepath.Join(utils.MakePath("assets"), "header.html"),
		filepath.Join(utils.MakePath("assets"), "footer.html"),
		filepath.Join(utils.MakePath("assets"), "topnav.html"),
	)
	if err != nil {
		return err
	}

	gm := newGoldmark()
	for _, page := range pages {
		// Using goldmark, convert the markdown to HTML
		mdBuffer := bytes.Buffer{}
		if err := gm.Convert(page.Content, &mdBuffer); err != nil {
			return err
		}

		// Setup the page params for template execution
		pageParams := setupPageParams(
			c.Template.Params,
			mdBuffer.String(),
			c.Template.Params["title"].(string),
		)
		// Put the output file in the build dir that the ExecuteTemplate function expects
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

		// Do the templating against the base template but with the individual page params
		// Output will be written to file
		if err := tmpl.ExecuteTemplate(file, "base_page.html", pageParams); err != nil {
			return err
		}
	}

	return nil
}

func setupPageParams(params map[string]interface{}, mainContent string, title string) map[string]interface{} {
	pageParams := map[string]interface{}{}
	for k, v := range params {
		pageParams[k] = template.HTML(v.(string))
	}
	pageParams["main_content"] = template.HTML(mainContent)
	pageParams["title"] = title
	return pageParams
}

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"makeHrefs":    makeHrefs,
		"makeNavTitle": makeNavTitleFromHref,
	}
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

func getAssets(path string) ([]string, error) {
	var assets []string

	dir, err := os.Open(utils.MakePath(path))
	if err != nil {
		return assets, err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return assets, err
	}

	for i, v := range files {
		// We only want to treat files as assets here. If there's nested
		// directories containg more assets, then getAssets needs to get
		// called with that nested path to handle that case separately
		if !v.IsDir() {
			assets = append(assets, strings.Split(files[i].Name(), ".")[0])
		}
	}

	return assets, nil
}

func makeHrefs(path string) ([]string, error) {
	var hrefs []string

	assets, err := getAssets(path)
	if err != nil {
		return hrefs, err
	}

	sort.Strings(assets)
	for _, v := range assets {
		hrefs = append(hrefs, makeHref(v, path))
	}

	return hrefs, nil
}

func makeHref(assetName, originalPath string) string {
	pathSplit := strings.Split(originalPath, "/")
	hrefRoot := pathSplit[len(pathSplit)-1]

	return fmt.Sprintf("/%s/%s", hrefRoot, assetName)
}

// a recursive function that returns all the pages in the directory
// and all subdirectories. The path is a parameter because this
// function is called recursively to get all the pages in the site underneath
// "assets/pages". The first call of this function should be with the empty string
// as path which represents the root of assets/pages
func getAssetPages(path string, params map[string]interface{}) ([]page.Page, error) {
	baseAssetPath := utils.MakePath("assets/pages")
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
			subPages, err := getAssetPages(filepath.Join(fullAssetPath, file.Name()), params)
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

func makeNavTitleFromHref(assetHref string) string {
	pathSplit := strings.Split(assetHref, "/")
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(
		strings.Join(strings.Split(pathSplit[len(pathSplit)-1], "_"), " "),
	)
}
