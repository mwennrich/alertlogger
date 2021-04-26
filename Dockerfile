FROM golang:1.16 AS build
ENV GO111MODULE=on
ENV CGO_ENABLED=0


COPY . /app
WORKDIR /app

RUN go build alertlogger.go

FROM alpine:3.13

WORKDIR /
USER 65534

COPY --from=build /app .

ENTRYPOINT ./alertlogger
