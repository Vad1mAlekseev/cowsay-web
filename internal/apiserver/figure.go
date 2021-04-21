package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vad1malekseev/cowsay-web/internal/views"
)

// FigureHandler is providing figure by name using "text" query param.
// If "text" does not provide, he set to "random".
// "random" text handling with fortune CLI program.
// If figure name does not exists, it will return 404 status code. Usage:
//
// /eyes?text=hello+world! - returns the figure with custom text.
//
// /eyes?text=There+your+are!&mode=plain - returns a plain/text.
//
// /eyes?text=random - returns a random text using fortune.
func (s *APIServer) FigureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fig := vars["figure"]

	query := r.URL.Query()
	isPlainMode := query.Get("mode") == "plain"

	text := query.Get("text")
	if text == "" {
		text = "random"
	}

	figs, err := s.cowsay.List()
	if err != nil {
		s.logger.Errorf("error getting figure: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	exists, pos := checkFigure(figs, fig)
	if !exists {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	preformatted, err := s.cowsay.Make(fig, text)
	if err != nil {
		s.logger.Errorf("error creating figure: %v", err)

		return
	}

	if isPlainMode {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		_, _ = w.Write(preformatted)

		return
	}

	prev, next := getNextPrevFig(figs, pos)

	figView := views.NewFigure(fig, next, prev, string(preformatted))
	if err := figView.HTML(w); err != nil {
		s.logger.Errorln("error parsing HTML-template with figure:", fig)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func checkFigure(figs []string, figure string) (exists bool, pos int) {
	for idx, fig := range figs {
		if fig == figure {
			exists = true
			pos = idx

			break
		}
	}

	return
}

func getNextPrevFig(figs []string, currFigPos int) (string, string) {
	nextIdx := currFigPos + 1
	if nextIdx >= len(figs) {
		nextIdx = 0
	}

	prevIdx := currFigPos - 1
	if prevIdx < 0 {
		prevIdx = len(figs) - 1
	}

	next := figs[nextIdx]
	prev := figs[prevIdx]

	return prev, next
}
