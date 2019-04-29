[![Go Report Card](https://goreportcard.com/badge/github.com/blind-oracle/nginx-prometheus)](https://goreportcard.com/report/github.com/blind-oracle/nginx-prometheus)
[![Coverage Status](https://coveralls.io/repos/github/blind-oracle/nginx-prometheus/badge.svg?branch=master)](https://coveralls.io/github/blind-oracle/nginx-prometheus?branch=master)
[![Build Status](https://travis-ci.org/blind-oracle/nginx-prometheus.svg?branch=master)](https://travis-ci.org/blind-oracle/nginx-prometheus)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

# Nginx log parser and Prometheus exporter
This service parses incoming syslog messages from Nginx sent over UDP and converts them into Prometheus metrics exported through the built-in HTTP server.

* If the prefix-list is specified then per-URI statistics are generated.
Using this feature without URI-prefix list is dangerous because it leads to an unbounded memory usage. It's therefore enabled only with a limited prefix list.

    The URI prefix list is a plain-text file with a single prefix per line, e.g.
    ```
    /api/call1.json
    /api/call2.json
    /api/call3.json
    ...
    ```

    URIs received from nginx are stripped of any query parameters - only the part before '?' is used.

* It optionally supports country lookup for client IPs using MaxMind GeoIP database.

## Nginx configuration snippet
```
log_format collector '$remote_addr|$scheme|$host|$request_method|$server_protocol|$request_uri|$status|$request_time|$request_length|$bytes_sent';
access_log syslog:server=1.1.1.1:1514,tag=nginx collector;
```

## Getting it
* Get binary directly from [releases](https://github.com/blind-oracle/nginx-prometheus/releases), only linux-amd64 build currently available: 

* Simple build, with latest versions of all dependencies:
    ```
    go get github.com/blind-oracle/nginx-prometheus
    ```

* Repoducible build, using dep to get pinned dependencies:
    Get *dep* binary, e.g.:
    ```
    go get -u github.com/golang/dep/cmd/dep
    ```

    ```
    go get github.com/blind-oracle/nginx-prometheus
    cd $GOPATH/src/github.com/blind-oracle/nginx-prometheus
    $GOPATH/bin/dep ensure
    go build
    ```

## Usage example
```
# ./nginx-prometheus \
 -debug false \
 -listenSyslog 0.0.0.0:1514 \
 -listenHTTP 0.0.0.0:1514 \
 -geoipCountryDB /etc/nginx-prometheus/country.mmdb \
 -uriPrefixFile /etc/nginx-prometheus/uriPrefixes.txt
```

* **debug** - prints every syslog message received if true (defaults to false)
* **listenSyslog** - ip:port on which to listen for UDP Syslog messages (defaults to 0.0.0.0:1514)
* **listenHTTP** - ip:port on which to listen for incoming HTTP requests from Prometheus (defaults to 0.0.0.0:11080)
* **geoipCountryDB** - path to MaxMind country GeoIP database (optional)
* **uriPrefixFile** - path to a URI prefix list (optional)

## Grafana dashboard
Add Prometheus source to Grafana and import *grafana-dashboard.json* for most common graphs.
