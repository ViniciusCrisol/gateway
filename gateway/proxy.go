package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ReverseProxy(route Route, request *http.Request, response http.ResponseWriter) error {
	targetURL, err := route.GetTargetURL()
	if err != nil {
		return err
	}
	reverseProxy := httputil.ReverseProxy{Director: reverseProxyDirector(targetURL)}
	reverseProxy.ServeHTTP(response, request)
	return nil
}

// reverseProxyDirector returns a function that modifies the request to reverse
// proxy it successfully, rewriting the request URL and Host to match the target URL
func reverseProxyDirector(targetURL *url.URL) func(*http.Request) {
	return func(request *http.Request) {
		query := targetURL.RawQuery
		request.Host = targetURL.Host
		request.URL.Host = targetURL.Host
		request.URL.Path = targetURL.Path
		request.URL.Scheme = targetURL.Scheme
		if query == "" || request.URL.RawQuery == "" {
			request.URL.RawQuery = query + request.URL.RawQuery
		} else {
			request.URL.RawQuery = query + "&" + request.URL.RawQuery
		}
	}
}
