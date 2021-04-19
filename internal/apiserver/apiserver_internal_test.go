package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var defaultCfg *Config = &Config{
	BindAddr:          "test",
	StaticURLPrefix:   "/test/",
	LogLevel:          "debug",
	ConnectionTimeout: time.Minute,
}

type fakeCowsay struct {
}

func (fakeCowsay) List() ([]string, error) {
	return []string{"test-1", "test-2", "test-3"}, nil
}

func (fakeCowsay) Make(string, string) ([]byte, error) {
	return []byte("test\nfigure\n"), nil
}

func newTestServer() *ApiServer {
	s := New(defaultCfg)

	s.configureServer()
	s.cowsay = &fakeCowsay{}

	return s
}

func Test_ApiServer_HandleHomePage(t *testing.T) {
	s := newTestServer()
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "valid",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			s.HomeHandler(res, req)
			assert.Equal(t, tc.expectedCode, res.Code)
		})
	}
}

func Test_ApiServer_HandleFigurePage(t *testing.T) {
	s := newTestServer()
	testCases := []struct {
		name                string
		url                 interface{}
		expectedCode        int
		expectedContentType string
	}{
		{
			name:                "valid",
			url:                 "/test-1",
			expectedCode:        http.StatusOK,
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "non-existent figure",
			url:                 "/press-F",
			expectedCode:        http.StatusNotFound,
			expectedContentType: "",
		},
		{
			name:                "plain mode",
			url:                 "/test-1?mode=plain",
			expectedCode:        http.StatusOK,
			expectedContentType: "text/plain; charset=utf-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tc.url.(string), nil)
			s.server.Handler.ServeHTTP(res, req)
			s.FigureHandler(res, req)
			assert.Equal(t, tc.expectedCode, res.Code)
			contentType := res.Header().Get("content-type")
			assert.Equal(t, tc.expectedContentType, contentType)
		})
	}
}
