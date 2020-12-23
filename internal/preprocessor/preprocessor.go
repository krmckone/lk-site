package preprocessor

import (
	"bytes"
	"html/template"

	"github.com/krmckone/ksite/internal/config"
)

// Run executes preprocessor steps on md in environment p
func Run(md []byte, p config.Params) ([]byte, error) {
	return runTemplate(md, p)
}

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
