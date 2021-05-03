package views

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed templates/index.gohtml
var indexTemplate string //nolint:gochecknoglobals

type IndexView struct {
	Title   string
	Figures []string
}

// HTML wraps the figure to html.
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

// NewHome creates new home view by name.
func NewHome(figures []string) IndexView {
	return IndexView{"Cowsay Web!", figures}
}
