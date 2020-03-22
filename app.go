// FlaminGo is a simple HTTP REST framework written in Go.
// 
// Based on a https://github.com/nathan-fiscaletti/synful (a PHP
// framework named Synful), FlaminGo aims to be a more streamlined,
// fast and simple version of it's predecessor; stripping out features
// such as the ORM means the final binary can perform much faster,
// especially under load.
//
// It's important to keep in mind that FlaminGo is a high level
// abstraction of a lot of the built in HTTP logic provided by Go off
// the shelf. The main purpose of FlaminGo is to work as a supporting
// library to those already existing features, and to abstract some of
// their logic to make it easier to consume. FlaminGo makes a large
// effort to ensure that all of the lower level structures already
// provided by Go are always exposed through the layer of abstraction.
package flamingo

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "os"
    "sync"
    "time"
    "flag"

    "golang.org/x/time/rate"
)

type AppMode uint

// Modes in which the application can run. Configured in the Mode
// property of the App structure.
const (
    ModeHTTP  AppMode = 1 << iota
    ModeHTTPS
)

// Used for generating Client IDs from the specified request. This is
// used primarily for rate limiting so that you can determine a unique
// identifier for each client that is initiating requests.
type ClientIDFactory func(*http.Request) string

// App is the default object for a restful application.
type App struct {
    // The mode in which to run the web server.
    Mode            AppMode
    // The address to how the HTTP server on.
    Address         string
    // The routes to make available to the server.
    Routes          Routes
    // A list of all Middleware that is applied globally to all
    // Requests and Responses that pass through this App. 
    //
    // You can add Middleware to this route easily using the
    // App.AddMiddleware() function.
    Middleware      []Middleware
    // The default serializer to use for all requests and responses.
    Serializer      *Serializer
    // The global rate limit for all requests.
    GlobalLimit     *rate.Limiter
    // The per-client rate limit.
    ClientLimit     *rate.Limiter
    // Used to generate a unique client identifier.
    ClientIDFactory ClientIDFactory
    // The TLS address if running in ModeHTTPS.
    TLSAddress      string
    // The TLS Certificate File if running in ModeHTTPS.
    TLSCertFile     string
    // The TLS Key File if running in ModeHTTPS.
    TLSKeyFile      string
    // Print the ACCESS logs to console
    PrintAccess     bool
    clientLimits    map[string]*rate.Limiter
}

// NewAppFromCLI will generate a new App using the parameters passed
// to the command line. See `-h` for a full description of the
// availabel parameters.
func NewAppFromCLI() App {
    addressPtr := flag.String(
        "http", "", 
        "the address on which to run (only applies to HTTP)")
    tlsAddressPtr := flag.String(
        "https", "", 
        "the address on which to run (only applies to HTTPS)")
    tlsCertFilePtr := flag.String(
        "https-cert", "", "the certificate file to use for HTTPS")
    tlsKeyFilePtr := flag.String(
        "https-key", "", "the key file to use for HTTPS")

    flag.Parse()

    mode := ModeHTTP
    if *addressPtr == "" && *tlsAddressPtr != "" {
        mode = ModeHTTPS
    } else if *addressPtr != "" && *tlsAddressPtr != "" {
        mode = mode | ModeHTTPS
    } else if *addressPtr == "" && *tlsAddressPtr == "" {
        fmt.Fprint(os.Stderr, 
            "error: missing -http or -https, see -h for help.\n")
        os.Exit(1)
    }

    return App {
        Mode: mode,
        Address: *addressPtr,
        TLSAddress: *tlsAddressPtr,
        TLSCertFile: *tlsCertFilePtr,
        TLSKeyFile: *tlsKeyFilePtr,
    }
}

// Listen will start listening for HTTP and HTTPS requests sent to the
// application and process them respectively.
func (app *App) Listen() {
    if len(app.Routes) < 1 {
        Log.Print("warning : no routes defined, configure in main.go")
    }

    for _,route := range app.Routes {
        Log.Printf(
            "initialize : loaded route %v %p\n",
            route.Path, route.Handler)
    }

    var wg sync.WaitGroup

    if ModeHTTP & app.Mode == ModeHTTP {
        wg.Add(1)
        go func() {
            Log.Printf(
                "initialize : http starting at %s\n", app.Address)
            log.Fatal(http.ListenAndServe(app.Address, app))
            wg.Done()
        }()
    }
    
    if ModeHTTPS & app.Mode == ModeHTTPS {
        wg.Add(1)
        go func() {
            Log.Printf(
                "initialize : https starting at %s\n", app.TLSAddress)
            log.Fatal(http.ListenAndServeTLS(
                app.TLSAddress, app.TLSCertFile, app.TLSKeyFile, app))
        }()
    }

    wg.Wait()
}

// ServeHTTP will handle the incoming request and respond to it.
// First, process any configured rate limits from GlobalLimit and
// ClientLimit. After rate limits have been processed, determine which
// route to pass the request to. If no route can be determined,
// respond with a default 404 Not Found. Otherwise, process any
// potential rate limits on the route and process the request.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    app.RateLimit(w, r)

    q := r.URL.RawQuery
    if q != "" {
        q = fmt.Sprintf("?" + q)
    }
    path := r.URL.Path[1:]

    for _,route := range app.Routes {
        if route.IsURL(path) && route.Method == r.Method {
            if route.Limit != nil && app.ClientIDFactory == nil {
                Log.Printf(
                    "%s : %s\n", "warning", 
                    "RouteLimit set but no ClientIDFactory")
            } else {
                if !route.Allowed(app.ClientIDFactory, r) {
                    w.WriteHeader(http.StatusTooManyRequests)
                    return
                }
            }

            serialized, contentType, response := 
            app.process(path, route, w, r)

            // Set the response headers
            for k,v := range response.Headers {
                w.Header().Set(k, v)
            }

            // Set the content type
            w.Header().Set("Content-Type", contentType)

            // Set the HTTP status code
            w.WriteHeader(response.HTTPStatus)

            // Output the response
            w.Write([]byte(serialized))

            if app.PrintAccess {
                Log.Printf(
                    "access %p %s %s%s handle %s %p result %v %v",
                    r, r.Method, path, q, route.Path, route.Handler,
                    response.HTTPStatus, time.Since(start))
            }
            return
        }
    }

    w.WriteHeader(http.StatusNotFound)
    if app.PrintAccess {
        Log.Printf(
            "access %p %s %s%s handle nil 0x0000000 result 404 %v",
            r, r.Method, path, q, time.Since(start))
    }
    return
}

// RateLimit will process any potentially configured rate limits for
// the specified request.
func (app *App) RateLimit(w http.ResponseWriter, r *http.Request) {
    if app.GlobalLimit != nil {
        if !app.GlobalLimit.Allow() {
            w.WriteHeader(http.StatusTooManyRequests)
            return
        }
    }

    if app.ClientLimit != nil {
        if app.ClientIDFactory != nil {
            clientid := app.ClientIDFactory(r)

            if app.clientLimits == nil {
                app.clientLimits = map[string]*rate.Limiter{}
            }

            if _,exists := app.clientLimits[clientid]; !exists {
                app.clientLimits[clientid] = rate.NewLimiter(
                    app.ClientLimit.Limit(), app.ClientLimit.Burst())
            }

            if !app.clientLimits[clientid].Allow() {
                w.WriteHeader(http.StatusTooManyRequests)
                return
            }
        } else {
            Log.Print(
              "warning : ClientLimit set but no ClientIDFactory\n")
        }
    }
}

// AddMiddleware adds the specified Middleware to the App.
func (app *App) AddMiddleware(mw Middleware) {
    app.Middleware = append(app.Middleware, mw)
}

// process will process the incoming data through the configured
// serializers, use the controller bound to the route to process the
// request and return the serialized response, the value for the
// content type header value, and the response object.
// If an error is encountered while serializing or deserializing the
// data 400 or 500 HTTP response code will be returned respectively.
func (app *App) process(path string, route *Route, 
                        w http.ResponseWriter, r *http.Request) (
                            string, string, *Response) {
    // Retrieve the body
    var body []byte
    var bodyErr error
    if r.ContentLength > 0 {
        body,bodyErr = ioutil.ReadAll(r.Body)
        if bodyErr != nil {
            serializer := DefaultSerializer
            if app.Serializer != nil {
                serializer = app.Serializer
            }
            serialized,_ := serializer.Serialize(
                map[string]interface{}{
                    "error":fmt.Sprintf("%v", bodyErr),
                },
            )
            contentType := serializer.ContentType
            
            return serialized, contentType, &Response{
                HTTPStatus: 500,
            }
        }
    }

    // Deserialize input data
    data := map[string]interface{}{}
    if len(body) > 0 {
        var deserr error
        if route.Serializer != nil {
            data, deserr = route.Serializer.Deserialize(string(body))
        } else if app.Serializer != nil {
            data, deserr = app.Serializer.Deserialize(string(body))
        } else {
            data, deserr = DefaultSerializer.Deserialize(string(body))
        }
        
        if deserr != nil {
            serializer := DefaultSerializer
            if app.Serializer != nil {
                serializer = app.Serializer
            }
            serialized,_ := serializer.Serialize(
                map[string]interface{}{
                    "error":fmt.Sprintf(
                        "failed to parse input data: %v", deserr),
                },
            )
            contentType := serializer.ContentType
            
            return serialized, contentType, &Response{
                HTTPStatus: 400,
            }
        }
    }

    // Construct the request
    request := Request {
        Path: path,
        Route: route,
        Data: data,
        Headers: RequestHeaders1D(r.Header),
        Params: RequestQuery1D(r.URL.Query()),
        HTTPRequest: r,
    }

    // Process any "before" middleware
    for _,mw := range app.Middleware {
        if mw.Before != nil {
            mw.Before(&request)
        }
    }

    response := route.Handle(&request)

    // Process any "after" middleware
    for _,mw := range app.Middleware {
        if mw.After != nil {
            mw.After(response)
        }
    }

    // Serialize the response
    var serialized string
    var err error
    var contentType string = ""
    if response.Serializer != nil {
        serialized,err = response.Serializer.Serialize(response.Data)
        contentType = response.Serializer.ContentType
    } else if route.Serializer != nil {
        serialized,err = route.Serializer.Serialize(response.Data)
        contentType = route.Serializer.ContentType
    } else if app.Serializer != nil {
        serialized,err = app.Serializer.Serialize(response.Data)
        contentType = app.Serializer.ContentType
    } else {
        serialized,err = DefaultSerializer.Serialize(response.Data)
        contentType = DefaultSerializer.ContentType
    }

    // Handle any serialization errors
    if err != nil {
        var lastser error
        serializer := DefaultSerializer
        if app.Serializer != nil {
            serializer = app.Serializer
        }
        serialized,lastser = serializer.Serialize(
            map[string]interface{}{
                "error":fmt.Sprintf(
                    "failed to serialize output data: %v", err),
            },
        )
        if lastser != nil {
            serialized = fmt.Sprintf(
                "failed to serialize output data: %v", err)
            contentType = "text/plain"
        } else {
            contentType = serializer.ContentType
        }
        response = &Response{
            HTTPStatus: 500,
        }
    }

    return serialized, contentType, response
}
