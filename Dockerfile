FROM golang:alpine as dev

RUN apk add git --no-cache

WORKDIR /app
COPY echoserver.go .
RUN go get github.com/gorilla/handlers


FROM dev as build
RUN  CGO_ENABLED=0 go build echoserver.go


FROM scratch
COPY --from=build /app/echoserver /echoserver
ENTRYPOINT ["/echoserver"]
EXPOSE 8090
