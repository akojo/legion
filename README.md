# Legion web server

`legion` is a small reverse proxy plus web server for serving static content.

`legion` can serve static files from local filesystem and reverse proxy requests
based on path, virtual host or both. It requires minimal configuration, which
makes it ideal for simple everyday use cases like serving content from current
directory or acting as a simple development server for a frontend + backend
application.

In addition `legion` acts as a showcase for several nifty features of Go's
standard HTTP library.

## Installing

If you have `go` installed:

```sh
go install github.com/akojo/legion@latest
```

Otherwise you can grab a binary for your architecture from
[Releases](https://github.com/akojo/legion/releases) page.

## Get Started

To serve files from current directory on port 8000:

```sh
legion
```

As an another example, to

- Start a server on port 3000, restricting access to local connections
- Serve files from `static/` under current directory
- Route requests to `/api/*` to a backend application over HTTP on localhost
  port 3001

```sh
legion -listen localhost:3000 -route /=static -route /api=http://localhost:3001
```

## Configuration

`legion` accepts configuration both from a [configuration
file](#configuration-file) and via [command-line
options](#command-line-options). If both are specified, command-line options
override simple values and append values to lists. For example, if you provide
routes both in a configuration file and via command-line, command-line routes
will be appended to the routes in configuration file.

### Routing

Most of `legion` configuration revolves around routes. Every route has two
parts: a *source* and a *target*. Together they define where incoming requests
will be routed.

#### Sources

Route source can be either

- A plain path, e.g. `/api` or `/`
- A path prefixed by a hostname, e.g. `www.example.com/api` or
  `www.example.com/`

Paths are prefix-matched to incoming request paths; longer paths always take
precedence over shorter ones. Prefixing a path with a hostname restricts a route
to match only when incoming request's `Host` header matches the given hostname.

Source paths are always absolute; omitting the leading slash will result in an
error.

#### Targets

Route target can be either

- A local filesystem path, e.g. `/var/www/html`. Paths can be relative, in which
  case they are interpreted relative to `legion`'s current working directory.
- An HTTP/HTTPS URL, e.g. `https://www.example.com/api/v1`

Given a local path `legion` serves files from the specified directory. If
incoming request specifies a directory and the target directory contains a file
named `index.html`, contents of `index.html` are returned instead. Otherwise
directory contents are served as an HTML page.

If requested filename is `index.html` and the file exists, request will be
redirected to its parent directory.

Given an HTTP/HTTPS URL `legion` acts as a reverse proxy, forwarding requests to
the specified address. `legion` adds usual [forwarding
headers](#forwarding-headers) to outgoing requests.

#### Path Rewriting

`legion` always performs path rewriting, stripping source path from incoming
requests and then appending the remainder to target path.

For example, given routes (see [configuration file syntax](#configuration-file)
below)

```yaml
static:
- source: /
  target: /var/www/html
proxy:
- source: /api
  target: http://localhost:3000/v1
```

incoming requests map to targets as follows:

- `/favicon.ico`: first route is selected and contents of file
  `/var/www/html/favicon.ico` are returned
- `/api/pets/1`: second route is selected and request is forwarded to
  `http://localhost:3000/v1/pets/1`

NB. in this case, according to `legion`'s routing rules, if file
`/var/www/html/index.html` exists request to `/` will return its contents.

### Default Configuration

When started with an empty configuration `legion` serves files from current
directory on port `8000` and writes access logs to stdout.

### Configuration File

Configuration file is written in YAML and allows specifying following settings

```yaml
listen: <addr>
loglevel: <info|warn|error>
tls:
  certificates:
  - <certificate1>
  - <certificate2>
  ...
routes:
  static:
  - <route1>
  - <route2>
  ...
  proxy:
  - <route1>
  - <route2>
  ...
```

| Name       | Description                                                                                               | Values                                    | Default |
|------------|-----------------------------------------------------------------------------------------------------------|-------------------------------------------|---------|
| `listen`   | Listen on given address (i.e. network interface) and port. Omitting address will listen on all interfaces | `host:port`, `ip:port` or `:port`         | `:8000` |
| `loglevel` | Set log minimum level. Request logs are suppressed when level is above `info`                             | `info\|warn\|error`                       | `info`  |

#### TLS

TLS section is optional. If defined it will contain a list of X.509
certificate key pair files. If multiple certificates are defined, an appropriate
one is selected based on incoming request.

```yaml
TLS:
  certificates:
  - certfile: <file>
    keyfile: <file>
```

| Name       | Description                                  |
|------------|----------------------------------------------|
| `certfile` | Certificate public key file in X.509 format  |
| `keyfile`  | Certificate private key file in X.509 format |

#### Routes

Routes define either static routes, serving files from local filesystem, or
proxy routes, relaying request to an upstream server.

```yaml
routes:
  static:
  - source: <path>|<hostname/path>
    target: <path>
  ...
  proxy:
  - source: <path>|<hostname/path>
    target: <url>
```

See [Routing](#routing) for more information on specifying routes.

### Command-line Options

In addition to configuration file, `legion` understands following command-line
flags. Command-line flags override values in configuration file.

- `-config <file>`

  Read configuration from specified file. Configuration file format is
  documented [above](#configuration-file)

- `-listen <address>`

  Listen on given address. Can be `hostname:port`, `ip:port` or just `:port`.

- `-loglevel info|warn|error`

  Set log level. Request logs are written with level "info" and can thus be
  suppressed by setting level to either "warn" or "error".

- `-route <source>=<target>`

  Route requests from `source` to `target`. See [Routing](#routing) for
  specifying sources and targets.

## Forwarding headers

When acting as a reverse proxy `legion` always adds forwarding headers
(`X-Forwarded-For`, `X-Forwarded-Host` and `X-Forwarded-Proto`) to outgoing
requests. If an existing `X-Forwarded-For` is found in inbound request it is
retained and client IP is appended to its value.

`Host` header of inbound requests is kept copied as-is to outgoing requests.

## Access log format

`legion` prints access logs in structured format using logfmt-compatible output.
An example request log line is

```text
time=2006-01-02T15:04:05Z07:00 level=INFO msg="200 GET /" method=GET proto=HTTP/1.1 path=/ address=localhost:8000 status=200 duration=591.8µs user_agent=curl/8.0.1
```
