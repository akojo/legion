# Legion web server

`legion` is a simple web server for serving static content, because sometimes
you just want to serve files.

## Usage

Invoked without command-line arguments `legion` serves files from current
directory on port `8000` and writes access logs to stdout.

`legion` understands following command-line flags:

- `-dir <directory>`

  Serve files under given directory. Defaults to current directory.

- `-listen <address>`

  Listen on given address. Can be `hostname:port` or just `:port`. Defaults to
  `:8000`.

- `-quiet`

  Disable request logging.

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
