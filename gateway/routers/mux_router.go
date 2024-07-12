package routers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ViniciusCrisol/gateway/gateway"
)

type MuxRouter struct {
	mux *http.ServeMux
}

func NewMuxRouter(mux *http.ServeMux) *MuxRouter {
	return &MuxRouter{
		mux: mux,
	}
}

func (muxRouter *MuxRouter) SetupRouter(routes []gateway.Route) {
	gateway.ValidateDependencies()
	for _, route := range routes {
		route.Validate()
		switch route.Method {
		case http.MethodGet:
			muxRouter.mux.HandleFunc(fmt.Sprintf("GET %s", route.Path), muxRouter.handleRequest(route))
		case http.MethodPut:
			muxRouter.mux.HandleFunc(fmt.Sprintf("PUT %s", route.Path), muxRouter.handleRequest(route))
		case http.MethodPost:
			muxRouter.mux.HandleFunc(fmt.Sprintf("POST %s", route.Path), muxRouter.handleRequest(route))
		case http.MethodPatch:
			muxRouter.mux.HandleFunc(fmt.Sprintf("PATCH %s", route.Path), muxRouter.handleRequest(route))
		case http.MethodDelete:
			muxRouter.mux.HandleFunc(fmt.Sprintf("DELETE %s", route.Path), muxRouter.handleRequest(route))
		default:
			panic("gateway method unknown: " + route.Method + " (available methods: GET, PUT, POST, PATCH and DELETE)")
		}
	}
}

func (muxRouter *MuxRouter) handleRequest(route gateway.Route) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		enrichedRoute := muxRouter.enrichRoute(route, request)

		if err := gateway.ApplyFilters(enrichedRoute, request, response); err != nil {
			gateway.HandleErr(err, route, request, response)
			return
		}
		if err := gateway.ReverseProxy(enrichedRoute, request, response); err != nil {
			gateway.HandleErr(err, route, request, response)
			return
		}
	}
}

func (muxRouter *MuxRouter) enrichRoute(route gateway.Route, request *http.Request) gateway.Route {
	enrichedRoute := muxRouter.enrichRouteWithPathParams(route, request)
	return enrichedRoute
}

func (muxRouter *MuxRouter) enrichRouteWithPathParams(route gateway.Route, request *http.Request) gateway.Route {
	enrichedRoute := route
	for _, part := range strings.Split(route.URL, "/") {
		if part != "" &&
			part[0] == ':' {
			param := part[1:]
			param = strings.Split(param, "/")[0]
			param = strings.Split(param, "?")[0]
			enrichedRoute.URL = strings.Replace(enrichedRoute.URL, ":"+param, request.PathValue(param), -1)
		}
	}
	return enrichedRoute
}
