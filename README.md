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

## Examples

To serve files from current directory on port 8000:

```sh
legion
```

As an another example, to

- Start a server on port 3000, restricting access to local connections
- Serve files from `static/` under current directory
- Route requests to `/api/*` to a backend application on local port 3001:

```sh
legion -listen localhost:3000 -route /=$(pwd)/static -route /api=http://localhost:3001
```

## Usage

Invoked without command-line arguments `legion` serves files from current
directory on port `8000` and writes access logs to stdout.

`legion` understands following command-line flags:

- `-route <source>=<target>`

  Route requests from `source` to `target`. Default value is `/=$PWD`, which
  serves the contents of current directory. For more information, see "Routing"
  below. Can be specified multiple times.

- `-listen <address>`

  Listen on given address. Can be `hostname:port` or just `:port`. Defaults to
  `:8000`.

- `-quiet`

  Disable request logging.

## Routing

`legion` routes are specified with `<source>=<target>`. The specification
consists of two parts: `source` which provides a routing rule match incoming
request, and `target` which tells where to fetch the response.

`source` can be:

- A plain path, e.g. `/api`
- A path prefixed by a hostname, e.g. `www.example.com/api`

In the latter case the path matches a request only when request's `Host` header
matches given hostname.

Source paths are always prefix rules and match a request for any path under the
given enpoint. Given several paths with a common prefix, longest matching path
is selected. You can read more about path matching from
[Go's ServeMux documentation](https://pkg.go.dev/net/http#ServeMux). `legion`
internally ensures that routes end with a `/` before installing handlers.

`target` can be:

- A local filesystem path, e.g. `/var/www/html`
- An HTTP/HTTPS URL, e.g. `https://www.example.com/api/v1`

In either case source path is first stripped from incoming requests and the
result is appended to target.

For example, given options

```text
 -route /=/var/www/html -route /api=https://www.example.com/v1
```

incoming paths map to actual requests:

- `/index.html` -> `/var/www/html/index.html`
- `/api/pets/1` -> `https://www.example.com/v1/pets/1`

## Access log format

`legion` prints access logs in machine-readable format. Fields are
space-separated and don't contain any embedded spaces save for the last part
that prints client's user agent string.

The fields are, in order (see [logger.go](./logger.go)):

- UTC Time in RFC3339 format
- Reponse status code
- Request Method
- Request Protocol
- Request Path, URL encoded
- Reponse length in bytes, followed by space and `B`
- Response time in milliseconds, with microsecond precision, followed by space
  and `ms`
- Rest of the line contains response user agent information (from request
  `User-Agent` header)

For example

```text
2022-01-23T19:44:34Z 200 GET HTTP/1.1 / 193 B 0.171 ms curl/7.77.0
```
