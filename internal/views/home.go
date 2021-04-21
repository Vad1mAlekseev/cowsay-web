package views

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed templates/index.html
var indexTemplate string

type IndexView struct {
	Title   string
	Figures []string
}

func (v *IndexView) HTML(dst io.Writer) error {
	indexTempl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return fmt.Errorf("error parsing home HTML template: %w", err)
	}

	if err := indexTempl.Execute(dst, v); err != nil {
		return fmt.Errorf("error executing home HTML template: %w", err)
	}

	return nil
}

func NewHome(figures []string) IndexView {
	return IndexView{"Cowsay Web!", figures}
}
