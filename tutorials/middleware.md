# Managing Middleware

Middleware acts as a form of pre-processing for Requests entering your Application and form of post-processing for Responses before they leave your Application. You can read more about the [`Middleware` structure](https://godoc.org/github.com/nathan-fiscaletti/galago#Middleware) in the official documentation.

## Overview

1. [Types of Middleware](#types-of-middleware)

## Types of Middleware

- **Before Middleware**

   Middleware that implement the [`Before` callback](https://godoc.org/github.com/nathan-fiscaletti/galago#Middleware.Before) will be executed immediately before a Routes Handler is called.

   ```go
    middleware := galago.Middleware {
        Before: func(response *galago.Request) {
            // implement the middleware
        },
    }
    ```

- **After Middleware**

   Middleware that implement the [`After` callback](https://godoc.org/github.com/nathan-fiscaletti/galago#Middleware.After) will be executed immediately after a Route handler returns a Response, but before that Response is sent to the user.

   ```go
    middleware := galago.Middleware {
        After: func(response *galago.Response) {
            // implement the middleware
        },
    }
    ```

- **Terminal Middleware**

   Middleware that implement the [`Terminate` callback](https://godoc.org/github.com/nathan-fiscaletti/galago#Middleware.Terminate) will be executed immediately after a Response is sent to the user.

   ```go
    middleware := galago.Middleware {
        Terminate: func(request *galago.Request, response *galago.Response) {
            // implement the middleware
        },
    }
    ```