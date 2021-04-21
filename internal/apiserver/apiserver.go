package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/vad1malekseev/cowsay-web/internal/cowsay"
)

type Metrics struct {
	totalRequests  *prometheus.CounterVec
	responseStatus *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
}

// The Cowsay CLI methods.
type Cowsay interface {
	List() ([]string, error)
	Make(string, string) ([]byte, error)
}

// The APIServer struct for serving cowsay-web.
type APIServer struct {
	server        *http.Server
	config        *Config
	logger        *logrus.Logger
	metrics       *Metrics
	prometheusCtx *prometheus.Registry

	cowsay Cowsay
}

// New creates a APIServer instance providing config.
func New(cfg *Config) *APIServer {
	registry := prometheus.NewRegistry()

	return &APIServer{
		config:        cfg,
		logger:        logrus.New(),
		prometheusCtx: registry,
		metrics: &Metrics{
			totalRequests: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "Number of get requests.",
				},
				[]string{"path"}),
			responseStatus: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "response_status",
					Help: "Status of HTTP response",
				},
				[]string{"status"},
			),
			httpDuration: promauto.With(registry).NewHistogramVec(prometheus.HistogramOpts{
				Name: "http_response_time_seconds",
				Help: "Duration of HTTP requests.",
			}, []string{"path"}),
		},
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

	_ = prometheus.Register(s.metrics.totalRequests)
	_ = prometheus.Register(s.metrics.responseStatus)
	_ = prometheus.Register(s.metrics.httpDuration)
	r.Use(s.prometheusMiddleware)

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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *APIServer) prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(s.metrics.httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		s.metrics.responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		s.metrics.totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}
