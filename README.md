# 🦩  FlaminGo

FlaminGo is a simple HTTP REST framework written in [Go](https://golang.org).

Based on a [PHP framework named Synful](https://github.com/nathan-fiscaletti/synful), FlaminGo aims to be a more streamlined, fast and simple version of it's predecessor; stripping out features such as the ORM means the final binary can perform much faster, especially under load. See [Performance](#performance) for more information.

It's important to keep in mind that FlaminGo is a high level abstraction of a lot of the built in HTTP logic provided by Go off the shelf. The main purpose of FlaminGo is to work as a supporting library to those already existing features, and to abstract some of their logic to make it easier to consume. FlaminGo makes a large effort to ensure that all of the lower level structures already provided by Go are always exposed through the layer of abstraction.

## Installation

```sh
$ go get github.com/nathan-fiscaletti/flamingo
```

```go
package main

import(
   "github.com/nathan-fiscaletti/flamingo"
)
```

## Demo & Documentation

Take a look at [the documentation](https://godoc.org/github.com/nathan-fiscaletti/flamingo/) and [the example file](./example/main.go) for further reading material and demonstrations of the framework.

Here is a simple Hello World route using the Framework.

_`hello_world.go`_
```go
package main

import (
    "https://github.com/nathan-fiscaletti/flamingo"
)

func main() {
    flamingo.NewAppFromCLI().AddController(
        flamingo.NewController().AddRoute(
            flamingo.NewRoute("GET", "hello/world", HelloWorld),
        ),
    ).Listen()
}

func HelloWorld(request flamingo.Request) *flamingo.Response {
    return flamingo.NewResponse(200, map[string]interface{} {
        "message": "Hello, World!",
    })
}
```

## Features

- **Written in Go**

  Writing the framework in Go allows for several performance enhancements that would not regularly be available. Namely, the ability to compile directly to machine code. See [Performance](#performance) for more information. FlaminGo attempts to maintain as much simplicity as possible while still providing the easiest implementation available.

- **HTTP / HTTPS**

   Both HTTP and HTTPS are supported by FlaminGo out of the box. See [Using TLS](#using-tls) for more information on using FlaminGo with HTTPS.

- **REST Components**

   FlaminGo supports several generic REST components that are not regularly available in the Go HTTP libraries. These include Middleware, Routing, Controllers, Serialization and Downloads. See the [example file](./example/main.go) for more information.
   
- **Consumable as a Library**

   FlaminGo does not require you to run it as a stand alone binary. Importing FlaminGo as a library into your existing HTTP project can be done quickly and easily to provide the same set of features available in FlaminGo to your existing web package.

## Why Go?

The older version of this framework was written in PHP. This was purely based on research around what languages were most commonly used for this type of framework at the time _(nearly five years ago)_. While using PHP did allow for some more in depth configuration and for features such as editing the code without re-compiling, coupled with the bloat and other annoyances of the older framework, it did not offer the raw performance you will get out of compiling your REST API down to machine code.

## What about [mod-go](https://github.com/idaunis/mod_go)?

In the future I hope to add support for mod-go, however the project has been abandoned for quite some time and I'm still on the fence about whether or not I will add support for it.

## What about an ORM?

If you'd like to use an ORM with FlaminGo to link in your database, I highly recommend looking into [gorm](https://github.com/jinzhu/gorm). It's a well maintained and highly recommended ORM for Go with over 17,000 stars. Since this library already exists and has a decent feature set and user base, I see no reason to write a new ORM for use with FlaminGo.

## Using TLS

Running FlaminGo with HTTPS is fairly easy. You'll first need to generate a certificate and a key to use. You can do this either using [certbot](https://certbot.eff.org/) or by generating one yourself.

To use certbot, run the following command (replacing `yourwebsite.com` with your domain name) and follow the on screen instructions to generate your key and certificate.

```sh
$ sudo certbot certonly --manual -d yourwebsite.com
```

Certbot will ask you to create a file located at a specific path on the website. You can do this by creating the following route using a `DownloadSerializer()`. Replace `<challenge_key>` with the last path element in the URL certbot requests, and replace `<challenge_data>` with the data expected to be found in the file.

```go
route := flamingo.NewRoute(
    "GET", ".well-known/acme-challenge/<challenge_key>",
    func(request flamingo.Request) *flamingo.Response {
        return flamingo.NewResponse(
            200, flamingo.DownloadSerializer().MakeRawData(
                "<challenge_data>",
            ),
        ).MakeDownload("<challenge_key>")
    },
)
```

You will then need to run the binary in HTTP mode for certbot to validate the request. If youa re using `flamingo.NewAppFromCLI()` you can do this by running the following command.

```sh
./myapp -http "yourwebsite.com:80"
```

Once you have your key and certificate, compile your binary. If you are using `flamingo.NewAppFromCLI()` you can use the following to run it with TLS. If you do not use `flamingo.NewAppFromCLI()`, simply configure the `TLSAddress`, `TLSCertFile` and `TLSKeyFile` properties of your `App`.

```sh
$ ./myapp -https "yourwebsite.com:443" -https-cert "./mycert.crt" -https-key "./mykey.key"
```

## License

FlaminGo is licensed under the Apache 2.0 License. See [LICENSE](./LICENSE) for more information.

## Performance

All tests run using a simple Hello World FlaminGo application and plain HTTP.

[**Baton**](https://github.com/americanexpress/baton)

```
$ ./baton -r 200000 -c 5 -u http://localhost:8080/hello/world
Configuring to send GET requests to: http://localhost:8080/hello/world
Generating the requests...
Finished generating the requests
Sending the requests to the server...
Finished sending the requests
Processing the results...


====================== Results ======================
Total requests:                                200000
Time taken to complete requests:          3.93532484s
Requests per second:                            50822
===================== Breakdown =====================
Number of connection errors:                        0
Number of 1xx responses:                            0
Number of 2xx responses:                       200000
Number of 3xx responses:                            0
Number of 4xx responses:                            0
Number of 5xx responses:                            0
=====================================================
```

[**Bombardier**](https://github.com/codesenberg/bombardier)

```
$ ./bombardier-darwin-amd64 -c 5 -n 200000 -l http://localhost:8080/hello/world
Bombarding http://localhost:8080/hello/world with 200000 request(s) using 5 connection(s)
 200000 / 200000 [==============================================================================================================================================================] 100.00% 42772/s 4s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec     44110.31    3090.28   49656.00
  Latency      111.53us    16.32us     2.35ms
  Latency Distribution
     50%   108.00us
     75%   119.00us
     90%   133.00us
     95%   147.00us
     99%   215.00us
  HTTP codes:
    1xx - 0, 2xx - 200000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:     8.75MB/s
```
