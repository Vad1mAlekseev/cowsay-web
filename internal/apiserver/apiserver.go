package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/vad1malekseev/cowsay-web/internal/cowsay"
)

type Cowsay interface {
	List() ([]string, error)
	Make(string, string) ([]byte, error)
}

type ApiServer struct {
	server *http.Server
	config *Config
	logger *logrus.Logger

	cowsay Cowsay
}

func New(cfg *Config) *ApiServer {
	return &ApiServer{config: cfg, logger: logrus.New()}
}

func (s *ApiServer) Run() error {
	if err := s.configureLogger(); err != nil {
		return fmt.Errorf("error configuring the logger: %v", err)
	}

	s.configureServer()
	s.configureCowsay()

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Fatalf("error trying to ListenAndServe: %v", err)
		}
	}()
	s.logger.Infoln("Server started listening on", s.config.BindAddr)

	return nil
}

func (s *ApiServer) WaitWithGracefulShutdown() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.ConnectionTimeout))
	defer cancel()
	_ = s.server.Shutdown(ctx)
	s.logger.Infoln("Shutting down...")
	os.Exit(0)
}

func (s *ApiServer) configureServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.HomeHandler)
	r.HandleFunc("/favicon.ico", s.FaviconHandler)
	r.HandleFunc("/{figure}", s.FigureHandler)

	staticHandler := http.StripPrefix(s.config.StaticURLPrefix, http.FileServer(http.Dir("web/static")))
	r.PathPrefix(s.config.StaticURLPrefix).Handler(staticHandler)

	timeout := time.Microsecond * time.Duration(s.config.ConnectionTimeout)
	s.server = &http.Server{
		Addr: s.config.BindAddr,
		// Avoid Slowloris attacks.
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
		IdleTimeout:  timeout,
		Handler:      r,
	}
}

func (s *ApiServer) configureLogger() error {
	lvl, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(lvl)
	return nil
}

func (s *ApiServer) configureCowsay() {
	s.cowsay = &cowsay.Cowsay{}
}
