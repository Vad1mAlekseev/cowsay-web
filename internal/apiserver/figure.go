package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vad1malekseev/cowsay-web/internal/views"
)

func (s *APIServer) FigureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := r.URL.Query()
	cow := vars["figure"]

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

	exists := false
	pos := 0

	for idx, fig := range figs {
		if fig == cow {
			exists = true
			pos = idx

			break
		}
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	preformatted, err := s.cowsay.Make(cow, text)
	if err != nil {
		s.logger.Errorf("error creating figure: %v", err)

		return
	}

	if isPlainMode {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		_, _ = w.Write(preformatted)

		return
	}

	nextIdx := pos + 1
	if nextIdx >= len(figs) {
		nextIdx = 0
	}

	prevIdx := pos - 1
	if prevIdx < 0 {
		prevIdx = len(figs) - 1
	}

	next := figs[nextIdx]
	prev := figs[prevIdx]

	figView := views.NewFigure(cow, next, prev, string(preformatted))
	if err := figView.HTML(w); err != nil {
		s.logger.Errorln("error parsing HTML-template with figure:", cow)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
