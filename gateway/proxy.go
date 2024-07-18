package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
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

			query := url.Values{}
			for key, vals := range route.TargetURL.Query() {
				query.Set(key, vals[0])
			}
			for key, vals := range request.URL.Query() {
				query.Set(key, vals[0])
			}
			request.URL.RawQuery = query.Encode()
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
