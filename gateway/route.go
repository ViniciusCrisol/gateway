package gateway

import "net/url"

type Route struct {
	URL     string
	Path    string
	Method  string
	Filters []RouteFilter
}

type RouteFilter struct {
	Name       string
	Properties map[string][]string
}

func (route *Route) Validate() {
	_, err := url.Parse(route.URL)
	if err != nil {
		panic("gateway route URL invalid: " + route.URL)
	}
}

func (route *Route) GetTargetURL() (*url.URL, error) {
	targetURL, err := url.Parse(route.URL)
	if err != nil {
		return nil, err
	}
	return targetURL, nil
}

func (route *Route) GetFilterProperties(filterName string) (map[string][]string, bool) {
	for _, filter := range route.Filters {
		if filter.Name == filterName {
			return filter.Properties, true
		}
	}
	return nil, false
}
