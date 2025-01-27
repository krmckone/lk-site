package templater

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/krmckone/lk-site/internal/config"
	"github.com/krmckone/lk-site/internal/utils"
	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Page holds data for templating a page
type Page struct {
	Title     string
	Content   []byte
	Template  []byte
	Params    map[string]interface{}
	AssetPath string
	BuildPath string
}

func (p *Page) String() string {
	return fmt.Sprintf(
		"Title: %s\nContent: %s\nTemplate: %s\nParams: %v\nAssetPath: %s\nBuildPath: %s",
		p.Title,
		string(p.Content),
		string(p.Template),
		p.Params,
		p.AssetPath,
		p.BuildPath,
	)
}

// BuildSite is for building the site. This includes templating HTML with markdown and
// putting images in the expected locations in the output
func BuildSite() error {
	dirs := []string{"build", "build/images", "build/js", "build/shaders"}
	for _, dir := range dirs {
		utils.Clean(dir)
		utils.Mkdir(dir)
	}
	assetDirs := []string{"images", "js", "shaders"}
	for _, dir := range assetDirs {
		utils.CopyAssetToBuild(dir)
	}

	gm := newGoldmark()

	c, err := config.ReadConfig("configs/config.yml")
	if err != nil {
		return err
	}

	pages, err := getPages("assets/pages", c.Template.Params)
	if err != nil {
		return err
	}

	funcs := map[string]interface{}{"makeHrefs": makeHrefs, "makeNavTitle": makeNavTitleFromHref}
	tmpl := template.New("base_page.html")
	tmpl, err = tmpl.Funcs(funcs).ParseFiles("assets/base_page.html", "assets/header.html", "assets/footer.html", "assets/topnav.html")
	if err != nil {
		return err
	}

	for _, page := range pages {
		// Using goldmark, convert the markdown to HTML
		mdBuffer := bytes.Buffer{}
		if err := gm.Convert(page.Content, &mdBuffer); err != nil {
			return err
		}

		// Setup the page params for template execution
		pageParams := make(map[string]interface{})
		for k, v := range c.Template.Params {
			if k == "githubIcon" || k == "linkedinIcon" {
				pageParams[k] = template.HTML(v.(string))
			} else {
				pageParams[k] = v
			}
		}
		pageParams["main_content"] = template.HTML(mdBuffer.String())
		pageParams["title"] = c.Template.Params["title"]
		// Put the output file in the build dir that the ExecuteTemplate function expects
		outputPath := filepath.Join(page.BuildPath, fmt.Sprintf("%s.html", page.Title))
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}
		file, err := os.Create(outputPath)
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
		),
	)
}

func getAssets(path string) ([]string, error) {
	var assets []string

	dir, err := os.Open(path)
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

func getPages(path string, params map[string]interface{}) ([]Page, error) {
	pages := []Page{}

	files, err := os.ReadDir(path)
	if err != nil {
		return pages, err
	}

	for _, file := range files {
		if file.IsDir() {
			subPages, err := getPages(filepath.Join(path, file.Name()), params)
			if err != nil {
				return pages, err
			}
			pages = append(pages, subPages...)
		}
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		// Read markdown content
		content, err := os.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return pages, err
		}

		buildPath := strings.ReplaceAll(path, "assets/pages", "build")
		// Create page with the markdown content
		title := strings.TrimSuffix(file.Name(), ".md")
		page := Page{
			Title:     title,
			Content:   content,
			Params:    params,
			AssetPath: path,
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
