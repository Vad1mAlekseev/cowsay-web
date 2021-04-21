package apiserver

import (
	"net/http"

	"github.com/vad1malekseev/cowsay-web/internal/views"
)

func (s *APIServer) HomeHandler(w http.ResponseWriter, r *http.Request) {
	figs, err := s.cowsay.List()
	if err != nil {
		s.logger.Errorf("error getting figures: %v", err)

		return
	}

	HomeFiller := views.NewHome(figs)

	if err := HomeFiller.HTML(w); err != nil {
		s.logger.Errorf("error while handling index page: %v", err)

		return
	}
}
