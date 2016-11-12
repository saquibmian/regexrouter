package regexrouter

import (
	"context"
	"net/http"
	"regexp"
	"fmt"
)

var anyHTTPMethod = "*"

// NewRouter creates a router using the supplied routes
func NewRouter() *RegexRouter {
	return &RegexRouter{
		routes: &[]route{},
	}
}

type route struct {
	Pattern  *regexp.Regexp
	Handlers map[string]http.Handler
}

// RegexRouter is a HTTP router that uses regexes to route to the correct
// http.HandlerFunc. The regexes support named capture groups, which are attached
// to http.Request.Context().
type RegexRouter struct {
	routes *[]route
}

// Handle registers a set of http.HandlerFuncs for the route provided, for the HTTP methods provided.
func (r RegexRouter) Handle(regex string, handlers map[string]http.Handler) {
	*r.routes = append(*r.routes, route{regexp.MustCompile(regex), handlers})
}

// HandleAny registers a http.HandlerFunc for the route provided for any HTTP method.
func (r RegexRouter) HandleAny(regex string, handler http.Handler) {
	r.Handle(regex, map[string]http.Handler{
		anyHTTPMethod: handler,
	})
}

func (r RegexRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	url := req.URL.Path

	for _, route := range *r.routes {
		if !route.Pattern.MatchString(url) {
			continue
		}

		handler := route.Handlers[req.Method]
		if handler == nil {
			handler = route.Handlers[anyHTTPMethod]
		}
		if handler == nil {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// extract URL parameters into the context
		urlParams := extractNamedCaptureGroups(route.Pattern, url)
		if len(urlParams) > 0 {
			ctx := req.Context()
			for name, value := range urlParams {
				ctx = context.WithValue(ctx, name, value)
			}
			req = req.WithContext(ctx)
		}
		fmt.Printf("regexrouter: extracted URL params: %q\n", urlParams)

		// invoke the handler
		handler.ServeHTTP(res, req)
		return
	}

	// no mapping
	http.NotFound(res, req)
}

func extractNamedCaptureGroups(regex *regexp.Regexp, pattern string) map[string]string {
	result := make(map[string]string)
	match := regex.FindStringSubmatch(pattern)
	if match == nil {
		return result
	}
	for i, name := range regex.SubexpNames() {
		if i == 0 {
			continue // first is the expression
		}
		result[name] = match[i+1]
	}
	return result
}
