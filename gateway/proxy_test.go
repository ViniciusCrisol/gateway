package gateway

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseProxy(t *testing.T) {
	var (
		expectedHost     = ""
		expectedPath     = "/product_kind"
		expectedStatus   = http.StatusCreated
		expectedRawQuery = "sort=-SKU&color=blue"
	)
	target := httptest.NewServer(
		http.HandlerFunc(
			func(response http.ResponseWriter, request *http.Request) {
				assert.Equal(t, expectedHost, request.Host)
				assert.Equal(t, expectedPath, request.URL.Path)
				assert.Equal(t, expectedRawQuery, request.URL.RawQuery)
				response.WriteHeader(expectedStatus)
			},
		),
	)
	expectedHost = strings.Replace(target.URL, "http://", "", -1)

	route := Route{
		URL: target.URL + "/product_kind?sort=-SKU",
	}
	request := http.Request{
		URL: &url.URL{RawQuery: "color=blue"},
	}
	response := httptest.NewRecorder()

	err := ReverseProxy(route, &request, response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, response.Code)
}

func TestReverseProxyWithAnInvalidURL(t *testing.T) {
	route := Route{
		URL: "https://www.ecommerce.com.br\\product_kind?sort=-SKU",
	}
	request := http.Request{
		URL: &url.URL{RawQuery: "color=blue"},
	}
	response := httptest.NewRecorder()

	err := ReverseProxy(route, &request, response)
	assert.Error(t, err)
}

func TestReverseProxyDirectorWithQueryParams(t *testing.T) {
	var (
		expectedHost     = "www.ecommerce.com.br"
		expectedPath     = "/product_kind"
		expectedScheme   = "https"
		expectedRawQuery = "sort=-SKU&color=blue"
	)
	request := http.Request{
		URL: &url.URL{RawQuery: "color=blue"},
	}
	targetURL, _ := url.Parse("https://www.ecommerce.com.br/product_kind?sort=-SKU")

	proxyDirector := reverseProxyDirector(targetURL)
	proxyDirector(&request)

	assert.Equal(t, expectedHost, request.Host)
	assert.Equal(t, expectedHost, request.URL.Host)
	assert.Equal(t, expectedPath, request.URL.Path)
	assert.Equal(t, expectedScheme, request.URL.Scheme)
	assert.Equal(t, expectedRawQuery, request.URL.RawQuery)
}

func TestReverseProxyDirectorWithNoQueryParams(t *testing.T) {
	var (
		expectedHost     = "www.ecommerce.com.br"
		expectedPath     = "/product_kind"
		expectedScheme   = "https"
		expectedRawQuery = "sort=-SKU"
	)
	request := http.Request{URL: &url.URL{}}
	targetURL, _ := url.Parse("https://www.ecommerce.com.br/product_kind?sort=-SKU")

	proxyDirector := reverseProxyDirector(targetURL)
	proxyDirector(&request)

	assert.Equal(t, expectedHost, request.Host)
	assert.Equal(t, expectedHost, request.URL.Host)
	assert.Equal(t, expectedPath, request.URL.Path)
	assert.Equal(t, expectedScheme, request.URL.Scheme)
	assert.Equal(t, expectedRawQuery, request.URL.RawQuery)
}
