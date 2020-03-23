## Using TLS

Running GalaGo with HTTPS is fairly easy. You can either configure the Galago binary directly with a TLS certificate and key, or if you choose to proxy it through Nginx or Apache, you can configure them as you normally would for TLS.

## Overview

1. [Using Apache or Nginx Proxy for TLS](#using-apache-or-nginx-proxy-for-tls)
2. [Configuring the Galago binary directly](#configuring-the-galago-binary-directly)
    - [Using Certbot to generate a Certificate and Key](#using-certbot-to-generate-a-certificate-and-a-key)
    - [Validating Certbot](#validating-certbot)
    - [Configure Galago](#configure-galago)

## Using Apache or Nginx Proxy for TLS

You can opt to use Apache or Nginx as a proxy for Galago and configure TLS there. Take a look at [Running Galago with Apache on Ubuntu](./apache.md) for more information.

## Configuring the Galago binary directly

It's recommended that you use [Certbot](https://certbot.eff.org/) to generate your TLS certificate and key. Certbot is a handy utility that takes most of the work out of configuring TLS.

### Using Certbot to generate a Certificate and Key

To use certbot, run the following command (replacing `yourwebsite.com` with your domain name) and follow the on screen instructions to generate your key and certificate.

```sh
$ sudo certbot certonly --manual -d yourwebsite.com
```

### Validating Certbot

Certbot will ask you to create a file located at a specific path on your website. You can do this by creating the following route using a `DownloadSerializer()`. Replace `<challenge_key>` with the last path element in the URL certbot requests, and replace `<challenge_data>` with the data expected to be found in the file.

```go
route := galago.NewRoute(
    http.MethodGet, ".well-known/acme-challenge/<challenge_key>",
    func(request galago.Request) *galago.Response {
        return galago.NewResponse(
            http.StatusOK, galago.DownloadSerializer().MakeRawData(
                "<challenge_data>",
            ),
        ).MakeDownload("<challenge_key>")
    },
)
```

Add this route to a `Controller` and add that `Controller` to your `App`. You will then need to run your Galago binary in HTTP mode for certbot to validate the request. If you are using `galago.NewAppFromCLI()` you can do this by running the following command.

```sh
./yourbinary -http "yourwebsite.com:80"
```

Once your application is listening, tell Certbot that it may proceed.

### Configure Galago

Once you have your key and certificate, compile your Galago binary. If you are using `galago.NewAppFromCLI()` you can use the following to run it with TLS. If you do not use `galago.NewAppFromCLI()`, simply configure the `TLSAddress`, `TLSCertFile` and `TLSKeyFile` properties of your `App`.

```sh
$ ./yourbinary -https "yourwebsite.com:443" -https-cert "./mycert.crt" -https-key "./mykey.key"
```

Alternately, you can choose to run your Galago binary as a service. See [Run Galago as a Service on Ubuntu](./service.md)