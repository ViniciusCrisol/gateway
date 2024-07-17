package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoute(t *testing.T) {
	t.Run("It should return a new Route", func(t *testing.T) {
		assert.NotPanics(t, func() {
			expectedPath := "/uuid"
			expectedMethod := http.MethodGet
			expectedTargetURL := "https://httpbin.org/uuid"

			route := NewRoute(
				expectedPath,
				expectedMethod,
				expectedTargetURL,
			)

			assert.Equal(t, expectedPath, route.Path)
			assert.Equal(t, expectedMethod, route.Method)
			assert.Equal(t, expectedTargetURL, route.TargetURL.String())
		})
	})

	// t.Run("It should panic when the path is invalid", func(t *testing.T) {
	// 	assert.Panics(t, func() {
	// 		NewRoute(
	// 			"\\uuid",
	// 			http.MethodGet,
	// 			"https://httpbin.org/uuid",
	// 		)
	// 	})
	// })

	t.Run("It should panic when the method is invalid", func(t *testing.T) {
		assert.Panics(t, func() {
			NewRoute(
				"/uuid",
				"UNKNOWN",
				"https://httpbin.org/uuid",
			)
		})
	})

	t.Run("It should panic when the target URL is invalid", func(t *testing.T) {
		assert.Panics(t, func() {
			NewRoute(
				"/uuid",
				http.MethodGet,
				"https://httpbin.org\\uuid",
			)
		})
	})
}

func TestIsPathValid(t *testing.T) {}

func TestIsMethodValid(t *testing.T) {
	t.Run("It should return true for http.MethodGet", func(t *testing.T) {
		assert.True(t, IsMethodValid(http.MethodGet))
	})

	t.Run("It should return true for http.MethodPut", func(t *testing.T) {
		assert.True(t, IsMethodValid(http.MethodPut))
	})

	t.Run("It should return true for http.MethodPost", func(t *testing.T) {
		assert.True(t, IsMethodValid(http.MethodPost))
	})

	t.Run("It should return true for http.MethodPatch", func(t *testing.T) {
		assert.True(t, IsMethodValid(http.MethodPatch))
	})

	t.Run("It should return true for http.MethodDelete", func(t *testing.T) {
		assert.True(t, IsMethodValid(http.MethodDelete))
	})

	t.Run("It should return true for GET", func(t *testing.T) {
		assert.True(t, IsMethodValid("GET"))
	})

	t.Run("It should return true for PUT", func(t *testing.T) {
		assert.True(t, IsMethodValid("PUT"))
	})

	t.Run("It should return true for POST", func(t *testing.T) {
		assert.True(t, IsMethodValid("POST"))
	})

	t.Run("It should return true for PATCH", func(t *testing.T) {
		assert.True(t, IsMethodValid("PATCH"))
	})

	t.Run("It should return true for DELETE", func(t *testing.T) {
		assert.True(t, IsMethodValid("DELETE"))
	})

	t.Run("It should return false for unknown methods", func(t *testing.T) {
		assert.False(t, IsMethodValid("UNKNOWN"))
	})
}
