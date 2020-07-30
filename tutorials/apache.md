# Running GalaGo with Apache on Ubuntu

You can easily configure GalaGo to run with Apache if you so wish. To do this, you will need to follow the instructions laid out below. 

## Overview

1. [Create the service](#create-the-service)
3. [Create the Apache Virtual Host](#create-the-apache-virtual-host)
4. [Configure TLS](#configure-tls-optional) _(Optional)_

## Create the Service

First, you should configure your GalaGo binary to run as a service. Follow [these instructions](./service.md) to do so.

## Create the Apache Virtual Host

You can technically use GalaGo now by sending requests to `localhost:8080`, however we want to use it with Apache. To do this, we'll need to create a virtual host file in Apache.

Start by creating a configuration file in `/etc/apache2/sites-enabled/` called `galago.mydomain.com.conf`.

> We assume that your domain name is `galago.mydomain.com`, however you should use whatever domain you have set up.

In this file, add the following

```conf
<VirtualHost *:80>
    ServerName galago.mydomain.com
    ProxyPreserveHost On
    <Proxy *>
        Order allow,deny
        Allow from all
    </Proxy>
    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/
</VirtualHost>
```

Next, you will need to enable two Apache mods.

```bash
$ sudo a2enmod proxy
$ sudo a2enmod proxy_http
```

Make sure you restart the Apache service when you are done.

```bash
$ service apache2 restart
```

Finally, test to make sure that everything is working.

```sh
$ curl https://galago.mydomain.com/example/download
Hello, World!
```

## Configure TLS _(Optional)_

You can optionally configure TLS. If you are running through Apache, I would not recommend that you use the built in TLS support that GalaGo provides and would instead recommend that you use [certbot](https://certbot.eff.org/). Once you have your Virtual Host set up, just run certbot as you normally would and everything should work as intended.

```sh
$ certbot
# follow on screen instructions
```

> For more information about GalaGo's TLS capabilities, read [Using TLS](./tls.md).