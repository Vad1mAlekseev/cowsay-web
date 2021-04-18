package apiserver

import (
	"github.com/vad1malekseev/cowsay-web/internal/views"
	"net/http"
)

func (s *ApiServer) HomeHandler(w http.ResponseWriter, r *http.Request) {
	figs, err := s.cowsay.List()
	if err != nil {
		s.logger.Errorf("error getting figures: %v", err)
	}

	HomeFiller := views.NewHome(figs)

	if err := HomeFiller.HTML(w); err != nil {
		s.logger.Errorf("error while handling index page: %v", err)
		return
	}
}
