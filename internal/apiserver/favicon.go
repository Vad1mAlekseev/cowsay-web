package apiserver

import (
	"io/ioutil"
	"net/http"
)

// FaviconHandler is providing the cowsay-web icon.
func (s *APIServer) FaviconHandler(w http.ResponseWriter, _ *http.Request) {
	icon, err := ioutil.ReadFile("web/static/img.png")
	if err != nil {
		s.logger.Errorf("error while handling favicon: %v", err)

		return
	}

	if _, err := w.Write(icon); err != nil {
		s.logger.Errorf("error while handling favicon: %v\n", err)

		return
	}
}
