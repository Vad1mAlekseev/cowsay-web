package apiserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestCfg() *Config {
	return &Config{
		BindAddr:          "test",
		StaticURLPrefix:   "/test/",
		LogLevel:          "debug",
		ConnectionTimeout: time.Minute,
	}
}

type fakeCowsay struct{}

func (fakeCowsay) List() ([]string, error) {
	return []string{"test-1", "test-2", "test-3"}, nil
}

func (fakeCowsay) Make(string, string) ([]byte, error) {
	return []byte("test\nfigure\n"), nil
}

func newTestServer() *APIServer {
	s := New(getTestCfg())

	s.configureServer()
	s.cowsay = &fakeCowsay{}

	return s
}

func Test_APIServer_HandleHomePage(t *testing.T) {
	t.Parallel()

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

	ctx := context.Background()

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			s.HomeHandler(res, req)
			assert.Equal(t, tc.expectedCode, res.Code)
		})
	}
}

func Test_APIServer_HandleFigurePage(t *testing.T) {
	t.Parallel()

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

	ctx := context.Background()

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tc.url.(string), nil)
			s.server.Handler.ServeHTTP(res, req)
			s.FigureHandler(res, req)
			assert.Equal(t, tc.expectedCode, res.Code)
			contentType := res.Header().Get("content-type")
			assert.Equal(t, tc.expectedContentType, contentType)
		})
	}
}
