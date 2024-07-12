package routers

import (
	"net/http"
	"testing"

	"github.com/ViniciusCrisol/gateway/gateway"

	"github.com/stretchr/testify/assert"
)

func TestMuxRouter_enrichRoute(t *testing.T) {
	var (
		muxRouter   MuxRouter
		expectedURL = "https://www.ecommerce.com.br/product_kind?sort=-SKU"
	)
	route := gateway.Route{
		URL: "https://www.ecommerce.com.br/:path_param?sort=-SKU",
	}
	request := http.Request{}
	request.SetPathValue("path_param", "product_kind")
	assert.Equal(t, expectedURL, muxRouter.enrichRouteWithPathParams(route, &request).URL)
}

func TestMuxRouter_enrichRouteWithPathParams(t *testing.T) {
	var (
		muxRouter   MuxRouter
		expectedURL = "https://www.ecommerce.com.br/product_kind?sort=-SKU"
	)
	route := gateway.Route{
		URL: "https://www.ecommerce.com.br/:path_param?sort=-SKU",
	}
	request := http.Request{}
	request.SetPathValue("path_param", "product_kind")
	assert.Equal(t, expectedURL, muxRouter.enrichRouteWithPathParams(route, &request).URL)
}
