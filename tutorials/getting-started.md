# My first GalaGo Application

In this tutorial we will write a very simple GalaGo application. This application will listen for requests to the endpoint `/user/{id}` and return the name of the associated user.

## Preparation

1. Make sure you have installed GalaGo using the following command.

   `go get github.com/nathan-fiscaletti/galago`

2. Create a new file called `main.go` and include the GalaGo code.

   ```go
   package main

   import (
       "github.com/nathan-fiscaletti/galago"
   )
   ```

## Creating your Application

The first thing you will need to do is create an Application. You can either implement your own or create one using the default method. [NewAppFromCLI()](https://godoc.org/github.com/nathan-fiscaletti/galago#NewAppFromCLI) will create a default application using the command line arguments passed to the binary at run time.

Make sure you call [`Listen()`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.Listen) on your Application after it has been initialized. This will start listening for requests.

You can read more about Applications in the [Managing Applications](./apps.md) documentation page.

```go
package main

import (
   "github.com/nathan-fiscaletti/galago"
)

func main() {
    app := galago.NewAppFromCLI()

    app.Listen()
}
```
## Creating your Controller

Next, you will need to create your Controller. A Controller effectively is a collection of Routes with traits applied globally. You can create a new Controller using the [NewController()](https://godoc.org/github.com/nathan-fiscaletti/galago#NewController) function.

After you have created your Controller, make sure you add it to your Application. You can add as many Controllers as you'd like, but for this example we will only need the one.

```go
package main

import (
   "github.com/nathan-fiscaletti/galago"
)

func main() {
    app := galago.NewAppFromCLI()
    controller := galago.NewController()



    app.AddController(controller)
    app.Listen()
}
```

## Creating your Route

Next, you will need to create a Route. A Route is effectively a path and matching logic that will be executed when that specific path is requested.

You can create a new route using the [NewRoute()](https://godoc.org/github.com/nathan-fiscaletti/galago#NewRoute) function.

You can add the Route to your Controller using the [controller.AddRoute](https://godoc.org/github.com/nathan-fiscaletti/galago#Controller.AddRoute) function.

Routes require a function in order to execute logic. The function should follow the [RouteHandler](https://godoc.org/github.com/nathan-fiscaletti/galago#RouteHandler) signature.

Note that the path for our Route uses a path segment wrapped in `{}`. This means that it is a path variable that can be retrieved at run time. We can retrieve this within our logic using [request.GetField("id")](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.GetField).

```go
package main

import (
   "github.com/nathan-fiscaletti/galago"
)

func main() {
    app := galago.NewAppFromCLI()
    controller := galago.NewController()

    controller.AddRoute(galago.NewRoute(
        "GET", "user/{id}", func(request galago.Request) *galago.Response {
            // implement logic
        },
    ))

    app.AddController(controller)
    app.Listen()
}
```

## Implementing the Route Logic

You're now ready to implement the logic for your Route. For this route, if the User ID requested is `0`, we will return the name `root`. Otherwise, we will return the name `John Doe`.

We can retrieve the value of the User ID by calling [request.GetField("id")](https://godoc.org/github.com/nathan-fiscaletti/galago#Request.GetField).

Once you are ready to creat the Response, you can create it using the [NewResponse()](https://godoc.org/github.com/nathan-fiscaletti/galago#NewResponse) function.

```go
package main

import (
   "github.com/nathan-fiscaletti/galago"
)

func main() {
    app := galago.NewAppFromCLI()
    controller := galago.NewController()

    controller.AddRoute(galago.NewRoute(
        "GET", "user/{id}", func(request galago.Request) *galago.Response {
            var userName = "John Doe"

            userIdPntr := request.GetField("id")
            if userIdPntr != nil {
                if *userIdPntr == "0" {
                    userName = "root"
                }
            }

            return galago.NewResponse(200, map[string]interface{}{
                "name": userName,
            })
        },
    ))

    app.AddController(controller)
    app.Listen()
}
```

## Testing your Application

Your finished application should now function.

- Compile

   ```sh
   $ go build main.go -o test
   ```

- Run

   ```sh
   $ ./test -http :8080
   2020/05/09 19:12:33 initialize : loaded route user/{id} 0x12aa4e0
   2020/05/09 19:12:33 initialize : http starting at :8080
   ```

- Test

   ```sh
   $ curl http://localhost:8080/user/0
   {"name":"root"}
   $ curl http://localhost:8080/user/123
   {"name":"John Doe"}
   ```

## Congratulations!

Congratulations on writing your first GalaGo application! Check out some of the other documentation for further information on using GalaGo.