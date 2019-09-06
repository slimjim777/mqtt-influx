package main

import (
	"github.com/slimjim777/mqtt-influx/config"
	"github.com/slimjim777/mqtt-influx/mqtt"
	"log"
	"os"
	"time"
)

const tickIntervalMins = 1

var mqttConn *mqtt.Connection

func main() {
	// Parse the env vars
	settings, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	defer mqttConn.Close()

	// Create/get the MQTT connection
	mqttConn, err = mqtt.GetConnection(settings)
	if err != nil {
		log.Printf("Error with MQTT connection: %v", err)
	}

	// On an interval...
	ticker := time.NewTicker(time.Minute * tickIntervalMins)
	for range ticker.C {
		// Create/get the MQTT connection
		mqttConn, err = mqtt.GetConnection(settings)
		if err != nil {
			log.Printf("Error with MQTT connection: %v", err)
			continue
		}
	}
	ticker.Stop()
}
