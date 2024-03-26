package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ViniciusCrisol/gateway/gateway"
	"github.com/ViniciusCrisol/gateway/gateway/routers"
)

var defaultErrorHandlerMessage = []byte(`{"message":"internal server error"}`)

func init() {
	gateway.RegisterDefaultErrorHandler(
		func(_ gateway.Route, _ *http.Request, response http.ResponseWriter) {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(defaultErrorHandlerMessage)
		},
	)
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	engine := gin.New()
	router := routers.NewGinRouter(engine)
	router.SetupRouter(
		[]gateway.Route{{URL: "https://httpbin.org/uuid", Path: "uuid", Method: http.MethodGet}},
	)
	server := http.Server{Addr: ":8080", Handler: engine}
	server.ListenAndServe()
}
