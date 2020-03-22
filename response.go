package flamingo

import (
    "fmt"
)

// Represents any Response that is passed through an App.
type Response struct {
    // The HTTP Status Code for the Response.
    HTTPStatus int
    // The Headers for the response. Easily set headers using the
    // response.SetHeader(key, val) function.
    Headers    map[string]string
    // The response data.
    Data       map[string]interface{}
    // The response Serializer. Easily set the Serializer using the
    // response.SetSerializer(serializer) function.
    Serializer *Serializer
}

// NewResponse creates a new response using the specified HTTP Status
// code and Data.
func NewResponse(status int, data map[string]interface{}) *Response {
    return &Response{
        HTTPStatus: status,
        Headers: map[string]string{},
        Data: data,
    }
}

// SetHeader sets the header specified in key to the value specified
// in val and returns the Response for further modification.
func (response *Response) SetHeader(key, val string) *Response {
    response.Headers[key] = val

    return response
}

// SetSerializer sets the serializer for the Response that will be
// used for serializing the data before returning it to the user. This
// function then returns the Response for further modification.
func (response *Response) SetSerializer(s *Serializer) *Response {
    response.Serializer = s
    return response
}

// MakeDownload will make this response function as a Download.
//
// It will set the Serializer for this Response to the 
// DownloadSerializer() serializer and it will set the
// Content-disposition header to attachment, and append the specified
// file name to it.
func (response *Response) MakeDownload(filename string) *Response {
    response.Serializer = DownloadSerializer()
    response.SetHeader(
        "Content-disposition",
        fmt.Sprintf("attachment; filename=%s", filename))

    return response
}