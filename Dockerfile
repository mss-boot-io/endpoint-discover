
FROM golang:alpine as builder

MAINTAINER lwnmengjing

#ENV GOPROXY https://goproxy.io/

WORKDIR /go/release
RUN apk update && apk add tzdata && apk add curl unzip procps ca-certificates

COPY go.mod ./go.mod
RUN go mod tidy
COPY . .
RUN pwd && ls

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o endpoint-discover .

FROM alpine

COPY --from=builder /go/release/endpoint-discover /