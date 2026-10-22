# Nginx log parser and Prometheus exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/blind-oracle/nginx-prometheus)](https://goreportcard.com/report/github.com/blind-oracle/nginx-prometheus)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

This service parses incoming syslog messages from Nginx sent over UDP or TCP and converts them into Prometheus metrics exported through the built-in HTTP server.

- If the prefix-list is specified then per-URI statistics are generated.
  Using this feature without URI-prefix list is dangerous because it leads to an unbounded memory usage. It's therefore enabled only with a limited prefix list.

  The URI prefix list is a plain-text file with a single prefix per line, e.g.

  ```bash
  /api/call1.json
  /api/call2.json
  /api/call3.json
  ```

  URIs received from nginx are stripped of any query parameters - only the part before '?' is used.

- It optionally supports country lookup for client IPs using MaxMind GeoIP database.

## Nginx configuration snippet

```nginxconf
log_format collector '$remote_addr|$scheme|$host|$request_method|$server_protocol|$request_uri|$status|$request_time|$request_length|$bytes_sent';
access_log syslog:server=1.1.1.1:1514,tag=nginx collector;
```

## Getting it

### Using Docker (recommended)

```bash
docker pull ghcr.io/blind-oracle/nginx-prometheus:latest

docker run -d \
  -p 1514:1514/udp \
  -p 1514:1514/tcp \
  -p 11080:11080 \
  -v /path/to/uriPrefixes.txt:/etc/nginx-prometheus/uriPrefixes.txt:ro \
  -v /path/to/country.mmdb:/etc/nginx-prometheus/country.mmdb:ro \
  ghcr.io/blind-oracle/nginx-prometheus:latest \
  -uriPrefixFile /etc/nginx-prometheus/uriPrefixes.txt \
  -geoipCountryDB /etc/nginx-prometheus/country.mmdb
```

Or use `docker-compose.yml` from the repository.

### Using RPM/DEB packages

Get RPM & DEB from [releases](https://github.com/blind-oracle/nginx-prometheus/releases), linux-amd64 only.

### Build from source

Requires Go 1.27.3+.

```bash
git clone https://github.com/blind-oracle/nginx-prometheus.git
cd nginx-prometheus
make build
```

## Usage example

If you're using RPM/DEB packages then options are set in `/etc/default/nginx-prometheus`

```bash
./nginx-prometheus \
  -debug false \
  -listenSyslog 0.0.0.0:1514 \
  -listenHTTP 0.0.0.0:11080 \
  -listenSyslogTCP 0.0.0.0:1514 \
  -geoipCountryDB /etc/nginx-prometheus/country.mmdb \
  -uriPrefixFile /etc/nginx-prometheus/uriPrefixes.txt
```

- **debug** — печатает каждое полученное syslog-сообщение (по умолчанию false)
- **listenSyslog** — ip:port для приёма UDP syslog-сообщений (по умолчанию 0.0.0.0:1514)
- **listenSyslogTCP** — ip:port для приёма TCP syslog-сообщений (по умолчанию отключено)
- **listenHTTP** — ip:port для HTTP-запросов от Prometheus (по умолчанию 0.0.0.0:11080)
- **geoipCountryDB** — путь к базе MaxMind GeoIP Country (опционально)
- **uriPrefixFile** — путь к файлу со списком URI-префиксов (опционально)

## Grafana dashboard

Add Prometheus source to Grafana and import _grafana-dashboard.json_ for most common graphs.
