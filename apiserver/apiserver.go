package apiserver

import (
	"cowsay-web/figures"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type ApiServer struct {
	handler http.Handler
}

func (s *ApiServer) SetupDefaultRoutes() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", IndexHandler)
	handler.HandleFunc("/figure", FigureHandler)
	handler.HandleFunc("/favicon.ico", FaviconHandler)
	handler.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	s.handler = handler
}

func (s ApiServer) ListenAndServe(addr string) error {
	if s.handler == nil {
		s.SetupDefaultRoutes()
	}

	err := http.ListenAndServe(addr, s.handler)
	if err != nil {
		return err
	}

	return nil
}

type IndexView struct {
	Title   string
	Figures []string
}

var indexTempl *template.Template
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if indexTempl == nil {
		idxTempl, err := ioutil.ReadFile(path.Join("templates", "index.html"))
		if err != nil {
			log.Printf("error while handling index page: %v", err)
			return
		}
		indexTempl, err = template.New("index").Parse(string(idxTempl))
		if err != nil {
			log.Printf("error while handling index page: %v", err)
			return
		}
	}


	availableFigrs, err := figures.Available()
	if err != nil {
		log.Printf("error while handling index page: %v", err)
		return
	}

	if err := indexTempl.Execute(w, IndexView{"Cowsay Web!", availableFigrs}); err != nil {
		log.Printf("error while handling index page: %v", err)
	}
}

func FigureHandler(w http.ResponseWriter, r *http.Request) {

}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	icon, err := ioutil.ReadFile(path.Join("static", "img.png"))
	if err != nil {
		log.Printf("error while handling favicon: %v\n", err)
		return
	}

	if _, err := w.Write(icon); err != nil {
		log.Printf("error while handling favicon: %v\n", err)
	}
}
