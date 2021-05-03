package views

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed templates/figure.gohtml
var figureTemplate string //nolint:gochecknoglobals

type FigureView struct {
	Title        string
	FigureName   string
	Preformatted string
	Next, Prev   string
}

// HTML wraps the figure to html.
func (v *FigureView) HTML(dst io.Writer) error {
	indexTempl, err := template.New("figure").Parse(figureTemplate)
	if err != nil {
		return fmt.Errorf("error parsing figure HTML template: %w", err)
	}

	if err := indexTempl.Execute(dst, v); err != nil {
		return fmt.Errorf("error executing figure HTML template: %w", err)
	}

	return nil
}

// NewFigure creates new figure view by name.
func NewFigure(figure, next, prev, preformatted string) *FigureView {
	return &FigureView{
		"Cowsay Web! | " + figure,
		figure,
		preformatted,
		next, prev,
	}
}
