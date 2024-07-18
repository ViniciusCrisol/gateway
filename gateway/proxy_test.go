package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseProxy(t *testing.T) {
	t.Run("It should handle a simple GET request", func(t *testing.T) {
		expectedPath := "/uuid"
		expectedMethod := http.MethodGet
		expectedResponse := "e345f3e6-ff6a-4a30-826b-0248619a9adf"

		server := httptest.NewServer(
			http.HandlerFunc(
				func(response http.ResponseWriter, request *http.Request) {
					assert.Equal(t, expectedPath, request.URL.Path)
					assert.Equal(t, expectedMethod, request.Method)
					response.Write([]byte(expectedResponse))
				},
			),
		)
		route := NewRoute(
			"/uuid",
			http.MethodGet,
			server.URL+expectedPath,
		)
		response := httptest.NewRecorder()

		ReverseProxy(
			route,
			httptest.NewRequest(route.Method, route.TargetURL.String(), nil),
			response,
		)

		assert.Equal(t, expectedResponse, response.Body.String())
	})

	t.Run("It should handle a GET request with target query params", func(t *testing.T) {
		expectedPath := "/names"
		expectedQuery := "limit=2"
		expectedMethod := http.MethodGet
		expectedResponse := "John\nMary"

		server := httptest.NewServer(
			http.HandlerFunc(
				func(response http.ResponseWriter, request *http.Request) {
					assert.Equal(t, expectedPath, request.URL.Path)
					assert.Equal(t, expectedQuery, request.URL.RawQuery)
					assert.Equal(t, expectedMethod, request.Method)
					response.Write([]byte(expectedResponse))
				},
			),
		)
		route := NewRoute(
			"/names",
			http.MethodGet,
			server.URL+expectedPath+"?"+expectedQuery,
		)
		response := httptest.NewRecorder()

		ReverseProxy(
			route,
			httptest.NewRequest(route.Method, route.TargetURL.String(), nil),
			response,
		)

		assert.Equal(t, expectedResponse, response.Body.String())
	})

	t.Run("It should handle a GET request overwriting the target query params", func(t *testing.T) {
		expectedPath := "/names"
		expectedQuery := "limit=4"
		expectedMethod := http.MethodGet
		expectedResponse := "John\nMary\nPeter\nChloe"

		server := httptest.NewServer(
			http.HandlerFunc(
				func(response http.ResponseWriter, request *http.Request) {
					assert.Equal(t, expectedPath, request.URL.Path)
					assert.Equal(t, expectedQuery, request.URL.RawQuery)
					assert.Equal(t, expectedMethod, request.Method)
					response.Write([]byte(expectedResponse))
				},
			),
		)
		route := NewRoute(
			"/names?limit=4",
			http.MethodGet,
			server.URL+expectedPath+"?"+expectedQuery,
		)
		response := httptest.NewRecorder()

		ReverseProxy(
			route,
			httptest.NewRequest(route.Method, route.TargetURL.String(), nil),
			response,
		)

		assert.Equal(t, expectedResponse, response.Body.String())
	})
}
