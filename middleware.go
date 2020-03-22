package flamingo

// Middleware acts as a form of pre processing for Requests and
// Responses. It's a useful tool for filtering requests.
//
// You can apply Middleware either globally to an App or on a per-
// route basis by directly applying it to the Route.
//
// When applied to either a specific Route or to an App, the Before
// function will be called on all Requests before they are passed
// to their respective Route. And the After funtion will be called on
// all Responses immediately after their RouteHandler has finished.
type Middleware struct {
    // Before is called Before any Request is handled.
    Before func(*Request)
    // After is called after a Response has been generated.
    After  func(*Response)
}
