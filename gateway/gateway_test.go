package gateway

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testableFilterName           = "TestableFilterName"
	testableFilterResponseStatus = http.StatusCreated

	testableErr                      = errors.New("testable error")
	testableErrHandlerResponseStatus = http.StatusForbidden

	testableDefaultErrorHandlerResponseStatus = http.StatusInternalServerError
)

type testableFilter struct{}

type testableErrHandler struct{}

func (_ *testableFilter) Name() string {
	return testableFilterName
}

func (_ *testableFilter) Apply(_ Route, _ *http.Request, response http.ResponseWriter) error {
	response.WriteHeader(testableFilterResponseStatus)
	return nil
}

func (_ *testableErrHandler) Error() error {
	return testableErr
}

func (_ *testableErrHandler) Handle(_ Route, _ *http.Request, response http.ResponseWriter) {
	response.WriteHeader(testableErrHandlerResponseStatus)
}

func testableDefaultErrorHandler(_ Route, _ *http.Request, response http.ResponseWriter) {
	response.WriteHeader(testableDefaultErrorHandlerResponseStatus)
}

func resetGatewayState() {
	filters = map[string]Filter{}
	filtersMt = sync.RWMutex{}

	customErrorHandlers = nil
	customErrorHandlersMt = sync.RWMutex{}

	defaultErrorHandler = func(route Route, request *http.Request, response http.ResponseWriter) {}
	defaultErrorHandlerMt = sync.RWMutex{}
	defaultErrorHandlerSet = false
}

func TestRegisterFilter(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()

	RegisterFilter(&testableFilter{})
	assert.Equal(t, 1, len(filters))
	assert.Equal(t, testableFilterName, filters[testableFilterName].Name())
}

func TestRegisterCustomErrorHandler(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()

	RegisterCustomErrorHandler(&testableErrHandler{})
	assert.Equal(t, 1, len(customErrorHandlers))
	assert.Equal(t, testableErr, customErrorHandlers[0].Error())
}

func TestRegisterDefaultErrorHandler(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()

	RegisterDefaultErrorHandler(
		func(_ Route, _ *http.Request, _ http.ResponseWriter) {},
	)
	assert.True(t, defaultErrorHandlerSet)
}

func TestValidateDependenciesWithDefaultErrorHandler(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()

	defer func() {
		assert.NotNil(t, recover())
	}()
	ValidateDependencies()
}

func TestValidateDependenciesWithoutDefaultErrorHandler(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()

	defaultErrorHandlerSet = true
	ValidateDependencies()
}

func TestApplyFilters(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()
	RegisterFilter(&testableFilter{})

	route := Route{
		Filters: []RouteFilter{{Name: testableFilterName}},
	}
	request := http.Request{}
	response := httptest.NewRecorder()

	err := ApplyFilters(route, &request, response)
	assert.NoError(t, err)
	assert.Equal(t, testableFilterResponseStatus, response.Code)
}

func TestApplyFiltersWithNonExistentFilter(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()

	route := Route{
		Filters: []RouteFilter{{Name: testableFilterName}},
	}
	request := http.Request{}
	response := httptest.NewRecorder()

	err := ApplyFilters(route, &request, response)
	assert.NoError(t, err)
}

func TestHandleErr(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()
	RegisterCustomErrorHandler(&testableErrHandler{})

	route := Route{
		Filters: []RouteFilter{{Name: testableFilterName}},
	}
	request := http.Request{}
	response := httptest.NewRecorder()

	HandleErr(testableErr, route, &request, response)
	assert.Equal(t, testableErrHandlerResponseStatus, response.Code)
}

func TestHandleErrWithDefaultHandler(t *testing.T) {
	resetGatewayState()
	defer resetGatewayState()
	RegisterDefaultErrorHandler(testableDefaultErrorHandler)

	route := Route{
		Filters: []RouteFilter{{Name: testableFilterName}},
	}
	request := http.Request{}
	response := httptest.NewRecorder()

	HandleErr(testableErr, route, &request, response)
	assert.Equal(t, testableDefaultErrorHandlerResponseStatus, response.Code)
}
