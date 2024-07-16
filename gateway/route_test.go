package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoute_ValidateWithValidURL(t *testing.T) {
	t.Run("It shouldn't panic when the given url is valid", func(t *testing.T) {
		route := Route{
			URL: "https://www.ecommerce.com.br/product_kind?sort=-SKU",
		}
		route.Validate()
	})

	t.Run("It should panic when the given url is invalid", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()
		route := Route{
			URL: "https://www.ecommerce.com.br\\product_kind?sort=-SKU",
		}
		route.Validate()
	})
}

func TestRoute_GetTargetURL(t *testing.T) {
	t.Run("It shouldn't return an error when the given url is valid", func(t *testing.T) {
		var (
			expectedHost     = "www.ecommerce.com.br"
			expectedPath     = "/product_kind"
			expectedScheme   = "https"
			expectedRawQuery = "sort=-SKU"
		)
		route := Route{
			URL: "https://www.ecommerce.com.br/product_kind?sort=-SKU",
		}

		targetURL, err := route.GetTargetURL()

		assert.NoError(t, err)
		assert.Equal(t, expectedHost, targetURL.Host)
		assert.Equal(t, expectedPath, targetURL.Path)
		assert.Equal(t, expectedScheme, targetURL.Scheme)
		assert.Equal(t, expectedRawQuery, targetURL.RawQuery)
	})

	t.Run("It should return an error when the given url is invalid", func(t *testing.T) {
		route := Route{
			URL: "https://www.ecommerce.com.br\\product_kind?sort=-SKU",
		}
		_, err := route.GetTargetURL()
		assert.Error(t, err)
	})
}

func TestRoute_GetFilterProperties(t *testing.T) {
	var (
		expectedFilter1Name = "1"
		expectedFilter2Name = "2"
		expectedFilter3Name = "3"
	)
	route := Route{
		Filters: []RouteFilter{
			{Name: "Filter1", Properties: map[string][]string{"FilterName": {expectedFilter1Name}}},
			{Name: "Filter2", Properties: map[string][]string{"FilterName": {expectedFilter2Name}}},
			{Name: "Filter3", Properties: map[string][]string{"FilterName": {expectedFilter3Name}}},
		},
	}
	_, ok := route.GetFilterProperties("Filter")
	assert.False(t, ok)

	filterName, ok := route.GetFilterProperties("Filter1")
	assert.True(t, ok)
	assert.Equal(t, expectedFilter1Name, filterName["FilterName"][0])

	filterName, ok = route.GetFilterProperties("Filter2")
	assert.True(t, ok)
	assert.Equal(t, expectedFilter2Name, filterName["FilterName"][0])

	filterName, ok = route.GetFilterProperties("Filter3")
	assert.True(t, ok)
	assert.Equal(t, expectedFilter3Name, filterName["FilterName"][0])
}

func TestDeleteCORSHeaders(t *testing.T) {
	response := http.Response{
		Header: http.Header{
			"Access-Control-Allow-Origin":      []string{"*"},
			"Access-Control-Allow-Methods":     []string{"DELETE, PATCH, POST, GET, PUT"},
			"Access-Control-Allow-Headers":     []string{"accept, origin, Content-Type, X-CSRF-Token, Authorization, Cache-Control, Content-Length, Accept-Encoding, X-Requested-With"},
			"Access-Control-Allow-Credentials": []string{"true"},
		},
	}

	deleteCORSHeaders(&response)

	assert.Equal(t, "", response.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "", response.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "", response.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "", response.Header.Get("Access-Control-Allow-Credentials"))
}
