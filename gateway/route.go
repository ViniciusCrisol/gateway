package gateway

import (
	"net/http"
	"net/url"
)

type Route struct {
	Path      string
	Method    string
	TargetURL url.URL
}

func NewRoute(
	path string,
	method string,
	rawTargetURL string,
) Route {
	if !IsPathValid(path) {
		panic("invalid path: " + path)
	}
	if !IsMethodValid(method) {
		panic("invalid method: " + method)
	}
	targetURL, err := url.Parse(rawTargetURL)
	if err != nil {
		panic("invalid target url: " + rawTargetURL)
	}
	return Route{
		Path:      path,
		Method:    method,
		TargetURL: *targetURL,
	}
}

func IsPathValid(path string) bool {
	return true
}

func IsMethodValid(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodPut ||
		method == http.MethodPost ||
		method == http.MethodPatch ||
		method == http.MethodDelete
}
