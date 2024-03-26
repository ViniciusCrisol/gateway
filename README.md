# Gateway

This project offers a versatile library for constructing an API Gateway using Golang routers. Drawing inspiration from
the [Spring Cloud Gateway](https://github.com/spring-cloud/spring-cloud-gateway), the library aims to simplify the
process of managing APIs. Currently, it integrates with Gin, a high-performance HTTP web framework for Golang.

## Getting Started

To integrate the Gateway project with your Golang application, follow the steps outlined below:

1. **Installation**: Begin by installing the Gateway library.
   You can either clone the repository or install it via a package manager like `go get`.

    ```bash
    go get github.com/ViniciusCrisol/gateway
    ```

2. **Configuration**: To configure your gateway, you just need to set a default error handler, that is used to return
   error responses to the end user. After that, connect the gateway to a Gin instance as shown in the following code:

   ```go
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
   ```

   Now, when a user sends a request to `localhost:8080/uuid`, the Gateway will initiate a reverse proxy
   to `https://httpbin.org/uuid`, allowing the user to interact with the resource located at that endpoint.
