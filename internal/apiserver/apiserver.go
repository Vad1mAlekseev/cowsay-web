package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/vad1malekseev/cowsay-web/internal/cowsay"
)

// The Cowsay CLI methods.
type Cowsay interface {
	List() ([]string, error)
	Make(string, string) ([]byte, error)
}

// The APIServer struct for serving cowsay-web.
type APIServer struct {
	server *http.Server
	config *Config
	logger *logrus.Logger

	cowsay Cowsay
}

// New creates a APIServer instance providing config.
func New(cfg *Config) *APIServer {
	return &APIServer{
		server: nil,
		config: cfg,
		logger: logrus.New(),
		cowsay: nil,
	}
}

// Run the server.
func (s *APIServer) Run() error {
	if err := s.configureLogger(); err != nil {
		return fmt.Errorf("error configuring the logger: %w", err)
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

// WaitWithShutdown checks for system interrupts and safe shutdown of the server if this happens.
func (s *APIServer) WaitWithShutdown() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ConnectionTimeout)
	_ = s.server.Shutdown(ctx)

	s.logger.Infoln("Shutting down...")
	cancel()
	os.Exit(0)
}

func (s *APIServer) configureServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.HomeHandler)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/favicon.ico", s.FaviconHandler)
	r.HandleFunc("/{figure}", s.FigureHandler)

	staticHandler := http.StripPrefix(s.config.StaticURLPrefix, http.FileServer(http.Dir("web/static")))
	r.PathPrefix(s.config.StaticURLPrefix).Handler(staticHandler)

	timeout := s.config.ConnectionTimeout
	s.server = &http.Server{
		Addr:        s.config.BindAddr,
		Handler:     r,
		TLSConfig:   nil,
		ReadTimeout: timeout,
		// Avoid Slowloris attacks.
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
	}
}

func (s *APIServer) configureLogger() error {
	lvl, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return fmt.Errorf("error parsing log level: %w", err)
	}

	s.logger.SetLevel(lvl)

	return nil
}

func (s *APIServer) configureCowsay() {
	s.cowsay = &cowsay.Cowsay{}
}
