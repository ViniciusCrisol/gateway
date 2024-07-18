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
}
