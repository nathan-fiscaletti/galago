package flamingo

import (
    "net/http"
    "net/url"
    "strings"
)

// Request represents any Request that passes through an App.
type Request struct {
    // The path specified in the Request.
    Path        string
    // The Route matching the Request.
    Route       *Route
    // The Data passed to the request with `-d`
    Data        map[string]interface{}
    // The Query Parameters passed to the Request.
    Params      map[string]string
    // The Headers available in the Request.
    Headers     map[string]string
    // The lower level http.Request structure.
    HTTPRequest *http.Request
}

// RequestQuery1D Converts a url.Values structure into a one
// dimensional map, using only the first value for each query
// parameter key.
func RequestQuery1D(query url.Values) map[string]string {
    res := map[string]string{}
    for k,q2d := range query {
        res[k] = q2d[0]
    }

    return res
}

// RequestHeaders1D Converts an http.Header structure into a one
// dimensional map, using only the first value for each header key.
func RequestHeaders1D(headers http.Header) map[string]string {
    res := map[string]string{}
    for k,h2d := range headers {
        res[k] = h2d[0]
    }

    return res
}

// GetField returns a pointer to the value for the Route Parameter
// matching the specified key. If none exists, nil is returned.
func (request *Request) GetField(key string) *string {
    pathProperties := request.Route.ParseProperties(request.Path)
    if val,exists := pathProperties[key]; exists {
        return &val
    }

    return nil
}

// GetHeader returns a pointer to the value for the Header matching
// the specified key. If none exists, nil is returned.
func (request *Request) GetHeader(key string) *string {
    key = http.CanonicalHeaderKey(key)
    if h,exists := request.Headers[key]; exists {
        return &h
    }

    return nil
}

// GetData returns the data found at the specified path. The path
// is delimeted with a period. For example, given the following data
//
// ```
// {"person":{"age":10}}
// ```
//
// You could acess the age of the person using the path `person.age`.
func (request *Request) GetData(path string) interface{} {
    var data interface{} = request.Data
    keys := strings.Split(path, ".")

    for _,key := range keys {
        m,isMap := data.(map[string]interface{})
        if isMap {
            if d,exists := m[key]; exists {
                data = d
            }
        } else {
            break
        }
        
    }

    return data
}

// GetQuery returns a pointer to the value for the Query Parameter
// matching the specified key. If none exists, nil is returned.
func (request *Request) GetQuery(param string) *string {
    if value,exists := request.Params[param]; exists {
        return &value
    }

    return nil
}
