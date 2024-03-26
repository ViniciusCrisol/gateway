package routers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ViniciusCrisol/gateway/gateway"
)

type GinRouter struct {
	engine *gin.Engine
}

func NewGinRouter(engine *gin.Engine) *GinRouter {
	return &GinRouter{
		engine: engine,
	}
}

func (ginRouter *GinRouter) SetupRouter(routes []gateway.Route) {
	gateway.ValidateDependencies()
	for _, route := range routes {
		route.Validate()
		switch route.Method {
		case http.MethodGet:
			ginRouter.engine.GET(route.Path, ginRouter.handleRequest(route))
		case http.MethodPut:
			ginRouter.engine.PUT(route.Path, ginRouter.handleRequest(route))
		case http.MethodPost:
			ginRouter.engine.POST(route.Path, ginRouter.handleRequest(route))
		case http.MethodPatch:
			ginRouter.engine.PATCH(route.Path, ginRouter.handleRequest(route))
		case http.MethodDelete:
			ginRouter.engine.DELETE(route.Path, ginRouter.handleRequest(route))
		default:
			panic("gateway method unknown: " + route.Method + " (available methods: GET, PUT, POST, PATCH and DELETE)")
		}
	}
}

func (ginRouter *GinRouter) handleRequest(route gateway.Route) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		enrichedRoute := ginRouter.enrichRoute(route, ctx)

		if err := gateway.ApplyFilters(enrichedRoute, ctx.Request, ctx.Writer); err != nil {
			gateway.HandleErr(err, route, ctx.Request, ctx.Writer)
			return
		}
		if err := gateway.ReverseProxy(enrichedRoute, ctx.Request, ctx.Writer); err != nil {
			gateway.HandleErr(err, route, ctx.Request, ctx.Writer)
			return
		}
	}
}

func (ginRouter *GinRouter) enrichRoute(route gateway.Route, ctx *gin.Context) gateway.Route {
	enrichedRoute := ginRouter.enrichRouteWithPathParams(route, ctx)
	return enrichedRoute
}

func (ginRouter *GinRouter) enrichRouteWithPathParams(route gateway.Route, ctx *gin.Context) gateway.Route {
	enrichedRoute := route
	for _, param := range ctx.Params {
		enrichedRoute.URL = strings.Replace(enrichedRoute.URL, ":"+param.Key, param.Value, -1)
	}
	return enrichedRoute
}
