# MQTT - Influx Data Logger

Quick-and-dirty MQTT data logger

Subscribes to one or more MQTT topics and writes raw message to an Influx database.
The data is expected in [influx data format](https://docs.influxdata.com/influxdb/v1.7/write_protocols/line_protocol_reference/) 
and no validation is done.

## Environment variables
```
INFLUXDB_URL      : default http://localhost:8086
INFLUXDB_USER     : username for Influxdb
INFLUXDB_PASSWORD : password for Influxdb
INFLUXDB_DATABASE : database for Influxdb

MQTT_HOST         : default localhost
MQTT_PORT         : default 1883
MQTT_USER         : username for the MQTT broker
MQTT_PASSWORD     : password for the MQTT broker
MQTT_TOPICS       : comma-separated list of topics to subscribe to
MQTT_TLS          : if TLS is needed - set to "true" or "false"
MQTT_TLS_CA       : path to the certificate authority file
MQTT_TLS_CERT     : path to the client certificate file
MQTT_TLS_KEY      : path to the client key file
```
The client ID defaults to the hostname of the server.
