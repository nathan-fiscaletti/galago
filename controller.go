package flamingo

// A Controller can group related HTTP request handling logic into a
// single structure.
type Controller struct {
	routes     []*Route
	middleware map[string][]Middleware
}

// NewController creates a new empty Controller.
func NewController() *Controller {
	return &Controller{
		routes:     []*Route{},
		middleware: map[string][]Middleware{},
	}
}

// GetAllRoutes retrieves all routes in the specified Controllers,
// applying any Controller level Middleware the the designated routes.
func GetAllRoutes(controllers ...*Controller) RouteCollection {
	res := RouteCollection{}

	for _, controller := range controllers {
		routes := controller.GetRoutes()
		for _, route := range routes {
			res = append(res, route)
		}
	}

	return res
}

// AddRoute adds a new Route to the Controller.
func (controller *Controller) AddRoute(route *Route) *Controller {
	controller.routes = append(controller.routes, route)
	return controller
}

// AddMiddleware adds a Middleware to the Controller and applies it to
// all Routes contained in this Controller.
func (controller *Controller) AddMiddleware(
	middleware Middleware,
) *Controller {
	if _, exists := controller.middleware["*"]; !exists {
		controller.middleware["*"] = []Middleware{}
	}

	controller.middleware["*"] = append(
		controller.middleware["*"], middleware)

	return controller
}

// AddMiddlewareFor adds a Middleware to the Controller and applies it
// to the routes specified in routes. Routes should be an array of
// paths used to identify the routes to apply the Middleware to.
func (controller *Controller) AddMiddlewareFor(
	routes []string, middleware Middleware,
) *Controller {
	for _, route := range routes {
		if _, exists := controller.middleware[route]; !exists {
			controller.middleware[route] = []Middleware{}
		}

		controller.middleware[route] = append(
			controller.middleware[route], middleware)
	}

	return controller
}

// GetRoutes retrieves all Routes within this Controller and applies
// the Middleware to the designated Routes.
func (controller *Controller) GetRoutes() RouteCollection {
	res := RouteCollection{}

	for _, route := range controller.routes {
		for path, mwc := range controller.middleware {
			if path == "*" || route.Path == path {
				for _, mw := range mwc {
					route.AddMiddleware(mw)
				}
			}
		}
		res = append(res, route)
	}

	return res
}
