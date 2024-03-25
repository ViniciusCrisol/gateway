package routers

import (
	"testing"

	"gateway/gateway"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinRouter_enrichRoute(t *testing.T) {
	var (
		ginRouter   GinRouter
		expectedURL = "https://www.ecommerce.com.br/product_kind?sort=-SKU"
	)
	route := gateway.Route{
		URL: "https://www.ecommerce.com.br/:path_param?sort=-SKU",
	}
	context := gin.Context{
		Params: gin.Params{{Key: "path_param", Value: "product_kind"}},
	}
	assert.Equal(t, expectedURL, ginRouter.enrichRoute(route, &context).URL)
}

func TestGinRouter_enrichRouteWithPathParams(t *testing.T) {
	var (
		ginRouter   GinRouter
		expectedURL = "https://www.ecommerce.com.br/product_kind?sort=-SKU"
	)
	route := gateway.Route{
		URL: "https://www.ecommerce.com.br/:path_param?sort=-SKU",
	}
	context := gin.Context{
		Params: gin.Params{{Key: "path_param", Value: "product_kind"}},
	}
	assert.Equal(t, expectedURL, ginRouter.enrichRouteWithPathParams(route, &context).URL)
}
