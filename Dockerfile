FROM golang:1.23.6-alpine AS build

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/legion

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/legion /
ENTRYPOINT ["/legion"]
