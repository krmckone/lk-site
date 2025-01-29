package page

import "fmt"

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
