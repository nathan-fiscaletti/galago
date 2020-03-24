# Managing Requests

Requests are used to represent the data that is received from your client. You can read more about the [`Request` structure](https://godoc.org/github.com/nathan-fiscaletti/galago#Request) in the official documentation.

## Overview

1. [Accessing Request Data](#accessing-request-data)
2. [Redirecting Requests](#redirecting-requests)
3. [Accessing the underlying HTTP Request](#accessing-the-underlying-http-request)

## Accessing Request Data

You can access different parts of the request, including Request Data, Query Parameters and Request Headers.

Each of these functions return a pointer to a string. If the key requested does not exist, or no data is found at the specified key, `nil` will be returned instead.

- **Request Data**

   To access Requests data, use the [`request.GetData(key)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.GetData).

   ```go
   value := request.GetData("name")
   ```

   > Request data is deserialized using the configured Serializer. See [Managing Serialization & Deserialization](./serialization.md) for more information.

- **Query Parameters**

   To access Query Parameters, use the [`request.GetQuery(key)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.GetQuery).

   ```go
   value := request.GetQuery("name")
   ```

- **Request Headers**

   To access Request Headers, use the [`request.GetHeader(key)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.GetHeader).

   ```go
   value := request.GetHeader("name")
   ```

## Redirecting Requests

You can redirect a request that comes into the framework using the [`request.Redirect(url, status)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.Redirect). This function takes a URL to which to redirect the Request and a Status to send back. It will return a Request object that represents the Redirect.

The URL can be either relative or absolute, the system is capable of handling both cases. Normally, this status will be one of `http.StatusPermanentRedirect`, `http.StatusTemporaryRedirect` or `status.MovedPermanently`.

```go
response := request.Redirect(
	"https://google.com/", http.StatusMovedPermanently,
)
```

## Accessing the underlying HTTP Request

You can access the underlying [`http.Request`](https://godoc.org/net/http#Request) using the [`HTTPRequest` property](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.HTTPRequest) of the `Request` structure.