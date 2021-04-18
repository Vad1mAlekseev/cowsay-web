package views

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed templates/figure.html
var figureTemplate string

type FigureView struct {
	Title        string
	FigureName   string
	Preformatted string
	Next, Prev   string
}

func (v *FigureView) HTML(dst io.Writer) error {
	indexTempl, err := template.New("figure").Parse(figureTemplate)
	if err != nil {
		return err
	}

	if err := indexTempl.Execute(dst, v); err != nil {
		return err
	}

	return nil
}

func NewFigure(figure, next, prev, preformatted string) *FigureView {
	return &FigureView{
		"Cowsay Web! | " + figure,
		figure,
		preformatted,
		next, prev,
	}
}
