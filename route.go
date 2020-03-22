package galago

import (
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/time/rate"
)

// Route represents a route within the application. Helps you route
// input to the appropriate action and respond with the proper data.
//
// You should not instantiate instances of this structure directly.
// Instead, use the NewRoute constructor function.
type Route struct {
	// The HTTP Method that should be applied to this Route.
	Method string
	// The path that this route should use. You can specify route
	// parameters using `{name}` format. For example, `users/{id}`.
	// Later, you can retrieve this value from the Request using
	// request.GetField("name").
	Path string
	// The regex match for the route. This is used to match incoming
	// URIs to the path of this Route.
	Match *regexp.Regexp
	// Processes a request that is sent to this Route.
	Handler RouteHandler
	// A list of all Middleware that is applied to this Route. You
	// can add Middleware to this route Easily using the
	// Route.AddMiddleware() function.
	Middleware []Middleware
	// The Serializer to use for Input / Output. See Serializer for
	// a more comprehensive description of what presedence serializers
	// take during a requests lifecycle.
	Serializer *Serializer
	// The Limiter applied to each client that requests this Route.
	// This Limiter is copied into a new Limiter for each client, so
	// it is only used as a reference for other Limiters.
	Limit        *rate.Limiter
	clientLimits map[string]*rate.Limiter
}

// RouteHandler handles Requests sent to a Route.
type RouteHandler func(Request) *Response

// RouteCollection represents a collection of Routes
type RouteCollection []*Route

// NewRoute creates a new Route with the specified HTTP method, path
// and RouteHandler.
func NewRoute(method, path string, handler RouteHandler) *Route {
	routeRegex := path

	var re = regexp.MustCompile(`(?m)\{([a-zA-Z0-9]+)\}`)
	for _, match := range re.FindAllStringSubmatch(path, -1) {
		routeRegex = strings.Replace(
			routeRegex, match[0], `(?P<`+match[1]+`>.+)`, -1)
	}

	return &Route{
		Method:     method,
		Path:       path,
		Match:      regexp.MustCompile(routeRegex),
		Handler:    handler,
		Middleware: []Middleware{},
	}
}

// Allowed determines if any rate limits are restricting the client
// specified in the result of the ClientIDFactory.
func (route *Route) Allowed(c ClientIDFactory, r *http.Request) bool {
	if route.Limit != nil {
		clientid := c(r)
		if route.clientLimits == nil {
			route.clientLimits = map[string]*rate.Limiter{}
		}

		if _, exists := route.clientLimits[clientid]; !exists {
			route.clientLimits[clientid] = rate.NewLimiter(
				route.Limit.Limit(), route.Limit.Burst())
		}

		return route.clientLimits[clientid].Allow()
	}

	return true
}

// AddMiddleware adds the specified Middleware to the Route.
func (route *Route) AddMiddleware(middleware Middleware) *Route {
	route.Middleware = append(route.Middleware, middleware)
	return route
}

// IsURL determines if the URL specified in url matches the Path
// set for this Route.
func (route *Route) IsURL(url string) bool {
	return route.Match.MatchString(url)
}

// ParseProperties parses the the route parameters found in a URL
// matching this Route. You should first determine if the URl matches
// this Route using IsUrl().
func (route *Route) ParseProperties(url string) map[string]string {
	match := route.Match.FindStringSubmatch(url)
	properties := make(map[string]string)
	for i, name := range route.Match.SubexpNames() {
		if i > 0 && i <= len(match) {
			properties[name] = match[i]
		}
	}

	return properties
}

// Handle will handle an incoming Request using this Route and return
// a Response.
//
// First all Before functions from the Middleware applied to this
// Route will be run on the Request, then the Request will be
// processed and a Response will be generated. That Response will then
// be run through all After functions from the Middleware applied to
// this Route and once completed, the Response will be returned.
func (route *Route) Handle(request *Request) *Response {
	// Process any "before" middleware
	for _, mw := range route.Middleware {
		if mw.Before != nil {
			mw.Before(request)
		}
	}

	// Process the request
	response := route.Handler(*request)

	// Process any "after" middleware
	for _, mw := range route.Middleware {
		if mw.After != nil {
			mw.After(response)
		}
	}

	return response
}
