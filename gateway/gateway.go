package gateway

import (
	"net/http"
	"sync"
)

type Filter interface {
	Name() string
	Apply(route Route, request *http.Request, response http.ResponseWriter) error
}

type ErrorHandler interface {
	Error() error
	Handle(route Route, request *http.Request, response http.ResponseWriter)
}

var (
	filters   = map[string]Filter{}
	filtersMt sync.RWMutex

	customErrorHandlers   []ErrorHandler
	customErrorHandlersMt sync.RWMutex

	defaultErrorHandler    func(route Route, request *http.Request, response http.ResponseWriter)
	defaultErrorHandlerMt  sync.RWMutex
	defaultErrorHandlerSet bool
)

func RegisterFilter(filter Filter) {
	filtersMt.Lock()
	defer filtersMt.Unlock()
	filters[filter.Name()] = filter
}

func RegisterCustomErrorHandler(customErrorHandler ErrorHandler) {
	customErrorHandlersMt.Lock()
	defer customErrorHandlersMt.Unlock()
	customErrorHandlers = append(customErrorHandlers, customErrorHandler)
}

func RegisterDefaultErrorHandler(errHandler func(route Route, request *http.Request, response http.ResponseWriter)) {
	defaultErrorHandlerMt.Lock()
	defer defaultErrorHandlerMt.Unlock()
	defaultErrorHandler = errHandler
	defaultErrorHandlerSet = true
}

func ValidateDependencies() {
	if !defaultErrorHandlerSet {
		panic("gateway default error handler not set")
	}
}

func ApplyFilters(route Route, request *http.Request, response http.ResponseWriter) error {
	filtersMt.RLock()
	defer filtersMt.RUnlock()

	for _, routeFilter := range route.Filters {
		filter, ok := filters[routeFilter.Name]
		if !ok {
			continue
		}
		if err := filter.Apply(route, request, response); err != nil {
			return err
		}
	}
	return nil
}

func HandleErr(err error, route Route, request *http.Request, response http.ResponseWriter) {
	customErrorHandlersMt.RLock()
	defer customErrorHandlersMt.RUnlock()
	defaultErrorHandlerMt.RLock()
	defer defaultErrorHandlerMt.RUnlock()

	for _, customErrorHandler := range customErrorHandlers {
		if customErrorHandler.Error() == err {
			customErrorHandler.Handle(route, request, response)
			return
		}
	}
	defaultErrorHandler(route, request, response)
}
