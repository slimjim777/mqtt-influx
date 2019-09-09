FROM golang:1.12 as builder1
COPY . ./src/github.com/slimjim777/mqtt-influx
WORKDIR /go/src/github.com/slimjim777/mqtt-influx
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/mqtt-influx -ldflags='-extldflags "-static"' cmd/mqtt-influx/main.go

# Copy the built applications to the docker image
FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder1 /go/bin/mqtt-influx /bin/mqtt-influx

ARG INFLUXDB_URL="http://localhost:8086"
ARG INFLUXDB_USER=""
ARG INFLUXDB_PASSWORD=""
ARG INFLUXDB_DATABASE="metrics"
ARG MQTT_HOST="localhost"
ARG MQTT_PORT="1883"
ARG MQTT_USER=""
ARG MQTT_PASSWORD=""
ARG MQTT_TOPICS="metrics,metrics/*"
ARG MQTT_TLS="false"
ARG MQTT_TLS_CA=""
ARG MQTT_TLS_CERT=""
ARG MQTT_TLS_KEY=""

ENV INFLUXDB_URL="${INFLUXDB_URL}"
ENV INFLUXDB_USER="${INFLUXDB_USER}"
ENV INFLUXDB_PASSWORD="${INFLUXDB_PASSWORD}"
ENV INFLUXDB_DATABASE="${INFLUXDB_DATABASE}"
ENV MQTT_HOST="${MQTT_HOST}"
ENV MQTT_PORT="${MQTT_PORT}"
ENV MQTT_USER="${MQTT_USER}"
ENV MQTT_PASSWORD="${MQTT_PASSWORD}"
ENV MQTT_TOPICS="${MQTT_TOPICS}"
ENV MQTT_TLS="${MQTT_TLS}"
ENV MQTT_TLS_CA="${MQTT_TLS_CA}"
ENV MQTT_TLS_CERT="${MQTT_TLS_CERT}"
ENV MQTT_TLS_KEY="${MQTT_TLS_KEY}"

ENTRYPOINT mqtt-influx