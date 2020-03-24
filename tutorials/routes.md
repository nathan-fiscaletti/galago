# Managing Routes

Routes allow you to route traffic that comes through your API to a specific function in your Application. You can read more about the [`Route` structure](https://godoc.org/github.com/nathan-fiscaletti/galago#Route) in the official documentation.

## Overview

1. [Creating a new Route](#creating-a-new-route)
    - [Managing Paths](#managing-paths)
2. [Applying Middleware to a Route](#applying-middleware-to-a-route)
3. [Using a custom Serializer with a Route](#using-a-custom-serializer-with-a-route)
4. [Applying a Rate Limit to a Route](#applying-a-rate-limit-to-a-route)
3. [Adding a Route to a Controller](#adding-a-route-to-a-controller)

## Creating a new Route

You should use the [`NewRoute()`](https://godoc.org/github.com/nathan-fiscaletti/galago#NewRoute) function when you want to create a new Route. Provide the function with a method, path and handler for the Route. This Route will only listen for requests to the specified path using the specified method.

```go
route := NewRoute(
    http.MethodGet, "my/path", 
    func(request galago.Request) *galago.Response {
        // return the response
    },
)
```

The function passed as the handler for the Route should follow the signature laid out in the [`RouteHandler` type definition](https://godoc.org/github.com/nathan-fiscaletti/galago#RouteHandler). For more information on Responses, see [Managing Responses](./responses.md).

### Managing Paths

Paths can contain segments called Route Parameters (or Fields) that can later be referenced within the handler for the Route. You can specify that a path segment is a field by wrapping it in `{}`. 

```go
"user/{id}"
```

You can also provide optional fields by wrapping the entire path segment (including it's preceding `/`) with `[]`.

```go
"user[/{id}]
```

To retrieve these fields in your Route handler, simply use the [`request.GetField(key)`](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.GetField) function. This function returns a pointer to a string. This pointer will bi `nil` if no value was found at the specified key.

```go
var user_id string
value := request.GetField("id")
if id != nil {
    user_id = *id
}
```

For more information on Requests, see [Managing Requests](./requests.md).

## Applying Middleware to a Route

You can apply Middleware to a route by using the [`route.AddMiddleware(middleware)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Route.AddMiddleware). This will apply the specified Middleware to any request that is handled by the Route.

```go
route.AddMiddleware(middleware)
```

## Using a custom Serializer with a Route

You can apply a Custom [`Serializer`](https://godoc.org/github.com/nathan-fiscaletti/galago#Serializer) to your Route using the [`Serializer`](https://godoc.org/github.com/nathan-fiscaletti/galago#Route.Serializer) property of your Route. This will force all requests that are sent to your Route to be parsable by the provided Serializer and format all Responses using the same Serializer. By default, Galago uses JSON for it's serialization and de-serialization. 

See [Managing Serialization and Deserialization](./serialization.md) for more information on Serializers.

## Applying a Rate Limit to a Route

> Galago uses [`rate.Limiter`](https://godoc.org/golang.org/x/time/rate#Limiter) for rate limiting. 

You can add a rate limit to a Route using the [`route.Limit`](https://godoc.org/github.com/nathan-fiscaletti/galago#Route.Limit) property.

When set, this will limit the rate at which requests can be sent to this Route from each individual client. It requires that you set the [`app.ClientIDFactory`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.ClientIDFactory) property of the Application that this Route belongs to in order to properly identify each client.

## Adding a Route to a Controller

Once you have prepared your Route, you can add it to a Controller using the [`controller.AddRoute(route)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Controller.AddRoute).

```go
controller.AddRoute(route)
```

For more information on Controllers, see [Managing Controllers](./controllers.md).