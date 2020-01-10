FROM golang:1.13-alpine AS build
ENV GO111MODULE=on
ENV CGO_ENABLED=0


COPY . /app
WORKDIR /app

RUN go build alertlogger.go

FROM alpine:latest

WORKDIR /

COPY --from=build /app .

ENTRYPOINT ./alertlogger
