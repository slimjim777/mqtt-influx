#!/usr/bin/env bash

export INFLUXDB_URL=http://localhost:8086
export INFLUXDB_USER=
export INFLUXDB_PASSWORD=
export INFLUXDB_DATABASE=metrics

export MQTT_HOST=localhost
export MQTT_PORT=1883
#export MQTT_USER=
#export MQTT_PASSWORD=
export MQTT_TOPICS=metrics,metrics/*
export MQTT_TLS=false
#export MQTT_TLS_CA=
#export MQTT_TLS_CERT=
#export MQTT_TLS_KEY=

go run cmd/mqtt-influx/main.go
