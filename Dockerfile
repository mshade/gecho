FROM golang:alpine as build

RUN apk add git --no-cache

WORKDIR /go/src
COPY echoserver.go .
RUN go get github.com/gorilla/handlers && \
  CGO_ENABLED=0 go build echoserver.go

FROM scratch
COPY --from=build /go/src/echoserver /echoserver
ENTRYPOINT ["/echoserver"]
EXPOSE 8090
