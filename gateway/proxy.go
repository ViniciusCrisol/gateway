package gateway

import (
	"net/http"
	"net/http/httputil"
)

func ReverseProxy(
	route Route,
	request *http.Request,
	response http.ResponseWriter,
) {
	reverseProxy := httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.Host = route.TargetURL.Host
			request.URL.Host = route.TargetURL.Host
			request.URL.Path = route.TargetURL.Path
			request.URL.Scheme = route.TargetURL.Scheme

			if request.URL.RawQuery != "" && route.TargetURL.RawQuery != "" {
				request.URL.RawQuery += "&" + route.TargetURL.RawQuery
				return
			}
			if request.URL.RawQuery == "" && route.TargetURL.RawQuery != "" {
				request.URL.RawQuery = route.TargetURL.RawQuery
				return
			}
		},
		ModifyResponse: func(response *http.Response) error {
			response.Header.Del("Access-Control-Max-Age")
			response.Header.Del("Access-Control-Allow-Origin")
			response.Header.Del("Access-Control-Allow-Methods")
			response.Header.Del("Access-Control-Allow-Headers")
			response.Header.Del("Access-Control-Expose-Headers")
			response.Header.Del("Access-Control-Allow-Credentials")
			return nil
		},
	}
	reverseProxy.ServeHTTP(response, request)
}
