# Managing Applications

Applications are used for managing your API and it's configuration. You can find the primary [`App` structure](https://godoc.org/github.com/nathan-fiscaletti/galago#App) in the official documentation.

## Overview

1. [Creating a new Application](#creating-a-new-appllication)
2. [Configuring your Application](#configuring-your-application)
   1. [Setting the Application Mode](#setting-the-application-mode)
   2. [Controllers](#controllers)
   3. [Middleware](#middleware)
   4. [Rate Limiting](#rate-limiting)
   5. [Custom Serializer](#custom-serializer)
   6. [Logging](#logging)
3. [Running your Application](#running-your-application)

## Creating a new Application

The easiest way to create a new App is to use the [`NewAppFromCLI()`](https://godoc.org/github.com/nathan-fiscaletti/galago#NewAppFromCLI) function. This will create a base `App` from the parameters passed in the command line. These include the following command line parameters.

```
  -http string
        the address on which to run (only applies to HTTP)
  -https string
        the address on which to run (only applies to HTTPS)
  -https-cert string
        the certificate file to use for HTTPS
  -https-key string
        the key file to use for HTTPS
```

```go
app := galago.NewAppFromCLI()
```

Alternately, you can create the App structure yourself.

```go
app := &galago.App{
    Mode: galago.ModeHTTP,
    Address: "localhost:80",
}
```

## Configuring your Application

There are several options available when configuring your application. They are described in detail here.

### Setting the Application Mode

The application mode ([`AppMode`](https://godoc.org/github.com/nathan-fiscaletti/galago#AppMode)) manages how this application should run. This tells the underlying system whether the binary should run as an HTTPS server, an HTTP server or both.

You can configure the `AppMode` of your application by setting `app.Mode`. You can provide multiple values by concatenating them with a `|` symbol. 

```go
app.Mode = galago.ModeHTTP | galago.ModeHTTPS
```

1. [`galago.ModeHTTP`](https://godoc.org/github.com/nathan-fiscaletti/galago#ModeHTTP)

   This will tell your Application to listen for regular HTTP requests at the address specified in `app.Address`.

2. [`galago.ModeHTTPS`](https://godoc.org/github.com/nathan-fiscaletti/galago#ModeHTTPS)

   This will tell your Application to listen for HTTPS requests at the address specified in `app.TLSAddress`. It is required that you set `app.TLSCertFile` and `app.TLSKeyFile` to the appropriate values for this to work properly. See [Configure Galago for TLS](./tls.md) for more information.

### Controllers

You can add a Controller to an Application using the [`app.AddController(controller)`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.AddController) function. Read more about Controllers in the [Managing Controllers](./controllers.md) tutorial.

### Middleware

You can add a Middleware to an Application using the [`app.AddMiddleware(middleware)`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.AddMiddleware) function. This Middleware will be applied to ALL requests to any routes within the Application. Read more about Middleware in the [Managing Middleware](./middleware.md) tutorial.

### Rate Limiting

> Galago uses [`rate.Limiter`](https://godoc.org/golang.org/x/time/rate#Limiter) for rate limiting. 

You can add a rate limit to an Application using the [`app.GlobalLimit`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.GlobalLimit) and the [`app.ClientLimit`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.ClientLimit) properties. 

- `app.GlobalLimit`

   This property is the rate limit to apply to the application itself. It limits the number of requests that can come through at any given time to avoid the system overloading.

- `app.ClientLimit`

   This property is the rate limit to apply to each client that consumes the API. It requires that you set the [`app.ClientIDFactory`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.ClientIDFactory) property of the Application in order to properly identify each client.

### Custom Serializer

You can apply a Custom [`Serializer`](https://godoc.org/github.com/nathan-fiscaletti/galago#Serializer) to your Application using the [`Serializer`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.Serializer) property of your Application. This will force all requests that are sent to your application to be parsable by the provided Serializer and format all Responses using the same Serializer. By default, Galago uses JSON for it's serialization and de-serialization. 

See [Managing Serialization and Deserialization](./serialization.md) for more information on Serializers.

### Logging

You can customize the logging for your application using the [`app.LogAccess`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.LogAccess) property. This will tell the Application whether or not it should be printing a log message for every request it receives.

## Running your Application

Once you have finished configuring your application, you can run it using the [`app.Listen`](https://godoc.org/github.com/nathan-fiscaletti/galago#App.Listen) function.

```go
func main() {
    app := NewAppFromCLI()

    // further configuration for the application

    app.Listen()
}
```