FROM golang:1.22 AS build
ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY . /app
WORKDIR /app

RUN go build alertlogger.go
RUN strip alertlogger

FROM scratch

WORKDIR /
USER 65534

COPY --from=build /app .

ENTRYPOINT ["./alertlogger"]
