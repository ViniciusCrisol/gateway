package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoute_ValidateWithValidURL(t *testing.T) {
	route := Route{
		URL: "https://www.ecommerce.com.br/product_kind?sort=-SKU",
	}
	route.Validate()
}

func TestRoute_ValidateWithInvalidURL(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()
	route := Route{
		URL: "https://www.ecommerce.com.br\\product_kind?sort=-SKU",
	}
	route.Validate()
}

func TestRoute_GetTargetURLWithValidURL(t *testing.T) {
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
}

func TestRoute_GetTargetURLWithInvalidURL(t *testing.T) {
	route := Route{
		URL: "https://www.ecommerce.com.br\\product_kind?sort=-SKU",
	}
	_, err := route.GetTargetURL()
	assert.Error(t, err)
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
