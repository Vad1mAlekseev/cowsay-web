package views

import (
	_ "embed"
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
		return err
	}

	if err := indexTempl.Execute(dst, v); err != nil {
		return err
	}

	return nil
}

func NewHome(figures []string) IndexView {
	return IndexView{"Cowsay Web!", figures}
}
