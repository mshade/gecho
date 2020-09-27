FROM golang:alpine as build

WORKDIR /app
COPY echoserver.go .
RUN CGO_ENABLED=0 go build echoserver.go

FROM scratch
COPY --from=build /app/echoserver /echoserver
ENTRYPOINT ["/echoserver"]
EXPOSE 8090
