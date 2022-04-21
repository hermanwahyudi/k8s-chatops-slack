# stage 1
FROM golang:1.12-alpine AS builder

RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates

RUN go get -d -v github.com/tkanos/gonfig

WORKDIR /app
COPY . /app

RUN GOOS=linux GOARCH=amd64 go build -o webhook 

# stage 2
FROM alpine:latest

WORKDIR /opt
COPY --from=builder /app /opt

EXPOSE 8081