FROM golang:1.20 AS build

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/legion

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/legion /
ENTRYPOINT ["/legion"]
