package apiserver

import (
	_ "embed"
	"io/ioutil"
	"net/http"
)

func (s *ApiServer) FaviconHandler(w http.ResponseWriter, r *http.Request) {
	icon, err := ioutil.ReadFile("web/static/img.png")
	if err != nil {
		s.logger.Errorf("error while handling favicon: %v", err)
		return
	}

	if _, err := w.Write(icon); err != nil {
		s.logger.Errorf("error while handling favicon: %v\n", err)
	}
}
