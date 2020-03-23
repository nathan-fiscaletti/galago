# Run Galago as a Service on Ubuntu

It's pretty easy to get your Galago binary running as a service on Ubuntu.

## Overview

1. [Create a new User](#create-a-new-user)
2. [Create the systemd service](#create-the-systemd-service)
3. [Start the service](#start-the-service)

## Create a new user

The service will need to run on a separate user in order to properly lock down the permissions for the application.

```sh
$ sudo useradd galago -s /sbin/nologin -M
```

## Create the systemd service

Create a file in the path `/lib/systemd/system/galago.service`

- Replace `/path/to/your/galago/binary` with the path to your Galago Binary
- Replace the `User` and `Group` with the user and group you created for running Galago

```ini
[Unit]
Description=Galago Framework
ConditionPathExists=/path/to/your/galago/binary
After=network.target
 
[Service]
Type=simple
User=galago
Group=galago
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/path/to/your/galago/
ExecStart=/path/to/your/galago/binary -http "localhost:8080"

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/galago
ExecStartPre=/bin/chown syslog:adm /var/log/galago
ExecStartPre=/bin/chmod 755 /var/log/galago
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=galago
 
[Install]
WantedBy=multi-user.target
```

## Start the service

```
$ sudo systemctl enable galago.service
$ sudo systemctl start galago
```

You can also re-start the service if need be using the following

```
$ sudo systemctl restart galago
```

You can view the logs from your service using the following

```
$ sudo journalctl -f -u galago
```