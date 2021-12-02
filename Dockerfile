FROM golang:1.17 AS build
ENV GO111MODULE=on
ENV CGO_ENABLED=0


COPY . /app
WORKDIR /app

RUN go build alertlogger.go
RUN strip alertlogger

FROM alpine:3.15

WORKDIR /
USER 65534

COPY --from=build /app .

ENTRYPOINT ./alertlogger
