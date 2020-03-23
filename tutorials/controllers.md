# Managing Controllers

The `Controller` effectively acts as a group of `Routes` and `Middleware` and is used for organizing these within your application.

You can read more about the [`Controller` structure](https://godoc.org/github.com/nathan-fiscaletti/galago#Controller), [`Routes`](https://godoc.org/github.com/nathan-fiscaletti/galago#Route) and [`Middleware`](https://godoc.org/github.com/nathan-fiscaletti/galago#Middleware) in the official documentation.



You can add multiple `Controllers` to your Application fairly easily. See [Managing Applications: Controllers](./apps.md#controllers) for more information.

## Overview

1. [Creating a new Controller](#creating-a-new-controller)
2. [Adding a Route to a Controller](#adding-a-route-to-a-controller)
3. [Applying Middleware to a Controller](#applying-middleware-to-a-controller)
4. [Adding a Controller to an Application](#adding-a-controller-to-an-application)

## Creating a new Controller

You should always create your controllers using the [`NewController()`](https://godoc.org/github.com/nathan-fiscaletti/galago#NewController) function.

```go
controller := NewController()
```

## Adding a Route to a Controller

The intention behind Controllers is to organize Routes. Adding a route to a controller is straight forward, simply call the [`controller.AddRoute(route)`](https://godoc.org/github.com/nathan-fiscaletti/galago#Controller.AddRoute) funciton to add the route to your controller. See [Managing Routes](./routes.md) for more information on creating routes.

```go
controller.AddRoute(route)
```

> Read more about Routes in the [Managing Routes](./routes.md) tutorial.

## Applying Middleware to a Controller

You can apply Middleware to a Controller in two ways. 

- [`controller.AddMiddleware(middleware)`](https://godoc.org/github.com/nathan-fiscaletti/galago#Controller.AddMiddleware)

   This will add the middleware to ALL Routes within the Controller; both existing ones and any routes you add in the future.

- [`controller.AddMiddlewareFor(routes, middleware)`](https://godoc.org/github.com/nathan-fiscaletti/galago#Controller.AddMiddlewareFor)

   This will add the middleware to only the routes specified in `routes`. `routes` should be an array of paths used to identify the routes.

   ```go
   controller.AddMiddlewareFor(
       []string{ "my/route", "other/route" },
       middleware,
   )
   ```

> You can read more about Middleware in the [Managing Middleware](./middleware.md) tutorial.

## Adding a Controller to an Application

You can add multiple Controllers to an Application, each one managing it's own set of Routes. Simply use the [`app.AddController(controller)`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.AddController) function to add a Controller to an Application.

```go
app.AddController(controller)
```