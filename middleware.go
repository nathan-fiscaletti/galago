package galago

// Middleware acts as a form of pre processing for Requests and
// Responses. It's a useful tool for filtering requests.
//
// You can apply Middleware either globally to an App , on a per-route
// basis by directly applying it to the Route, or to either all routes
// in a Controller or specific routes in a Controller.
type Middleware struct {
	// Before is called Before any Request is handled.
	Before func(*Request)
	// After is called after a Response has been generated.
	After func(*Response)
	// Terminate is called after the response has been sent. Request
	// can be nil under certain conditions.
	Terminate func(*Request, *Response)
}
