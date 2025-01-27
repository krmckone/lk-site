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
		pageParams["title"] = page.Title
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

/* func (p Page) exec(gm goldmark.Markdown) error {
	// Template the page content itself before letting goldmark convert from md to HTML
	tmpl, err := template.New("template").Parse(string(p.Content))
	if err != nil {
		return err
	}
	templBuffer := new(bytes.Buffer)
	if err = tmpl.Execute(templBuffer, p.Params); err != nil {
		return err
	}
	mdBuffer := new(bytes.Buffer)
	if err := gm.Convert(templBuffer.Bytes(), mdBuffer); err != nil {
		return err
	}
	p.Params["main_content"] = template.HTML(mdBuffer.String())

	tmpl, err = template.New("template").Parse(string(p.Template))
	if err != nil {
		return err
	}

	templBuffer = new(bytes.Buffer)
	if err = tmpl.Execute(templBuffer, p.Params); err != nil {
		return err
	}

	pathRoot := fmt.Sprintf("build/%s", p.BuildPathRoot)
	os.MkdirAll(pathRoot, 0700)
	utils.WriteFile(fmt.Sprintf("%s/%s.html", pathRoot, p.Title), templBuffer.Bytes())

	return nil
}

func runComponents(gm goldmark.Markdown, c *config.Config) error {
	// Do topnav content
	if err := runComponentTemplate(gm, c, "topnav"); err != nil {
		return err
	}

	// Do header content
	if err := runComponentTemplate(gm, c, "header"); err != nil {
		return err
	}

	// Do footer content
	if err := runComponentTemplate(gm, c, "footer"); err != nil {
		return err
	}

	// Steam Deck top 50
	if err := runComponentTemplate(gm, c, "steam_deck_top_50"); err != nil {
		return err
	}

	return nil
} */

/* func runComponentTemplate(gm goldmark.Markdown, c *config.Config, name string) error {
	buf := new(bytes.Buffer)
	md, err := utils.ReadFile(fmt.Sprintf("assets/components/%s.md", name))
	if err != nil {
		return err
	}

	switch name {
	case "topnav":
		md, err = runNavTemplate(md, c.Template.Params)
	case "steam_deck_top_50":
		md, err = runSteamDeckTop50Template(md, c.Template.Params)
	default:
		md, err = runTemplate(md, c.Template.Params)
	}
	if err != nil {
		return err
	}
	if err := gm.Convert(md, buf); err != nil {
		return err
	}
	c.Template.Params[name] = buf.String()

	return nil
} */

func runTemplate(md []byte, p config.Params) ([]byte, error) {
	tmpl, err := template.New("template").Parse(string(md))
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	if err = tmpl.Execute(buffer, p); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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

// func runNavTemplate(md []byte, p config.Params) ([]byte, error) {
// 	funcs := map[string]interface{}{"makeHrefs": makeHrefs, "makeNavTitle": makeNavTitleFromHref}
// 	tmpl, err := template.New("topnav").Funcs(funcs).Parse(string(md))
// 	if err != nil {
// 		return nil, err
// 	}

// 	buffer := new(bytes.Buffer)

// 	if err = tmpl.Execute(buffer, p); err != nil {
// 		return nil, err
// 	}

// 	return buffer.Bytes(), nil
// }

// TODO: Can we generalize component generation?
/* func runSteamDeckTop50Template(md []byte, p config.Params) ([]byte, error) {
	funcs := map[string]interface{}{"topFiftySteamDeckGames": steamapi.GetTopFiftySteamDeckGames}
	tmpl, err := template.New("topFiftySteamDeckGames").Funcs(funcs).Parse(string(md))
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)

	if err = tmpl.Execute(buffer, p); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
*/

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

func newPage(title string, content []byte, template []byte, params map[string]interface{}, path string) (Page, error) {
	return Page{
		Title:     title,
		Content:   content,
		Template:  template,
		Params:    params,
		AssetPath: path,
		BuildPath: strings.ReplaceAll(path, "assets", "build"),
	}, nil
}

func makeNavTitleFromHref(assetHref string) string {
	pathSplit := strings.Split(assetHref, "/")
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(
		strings.Join(strings.Split(pathSplit[len(pathSplit)-1], "_"), " "),
	)
}
