package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Default settings
const (
	defaultMQTTHost    = "localhost"
	defaultMQTTPort    = "1883"
	defaultClientID    = "mqtt_influx"
	defaultInfluxdbURL = "http://localhost:8086"

	envInfluxdbURL      = "INFLUXDB_URL"
	envInfluxdbUser     = "INFLUXDB_USER"
	envInfluxdbPassword = "INFLUXDB_PASSWORD"
	envInfluxdbDatabase = "INFLUXDB_DATABASE"

	envMQTTHost     = "MQTT_HOST"
	envMQTTPort     = "MQTT_PORT"
	envMQTTUser     = "MQTT_USER"
	envMQTTPassword = "MQTT_PASSWORD"
	envMQTTTopics   = "MQTT_TOPICS"
	envMQTTUseTLS   = "MQTT_TLS"
	envMQTTTLSCA    = "MQTT_TLS_CA"
	envMQTTTLSCert  = "MQTT_TLS_CERT"
	envMQTTTLSKEY   = "MQTT_TLS_KEY"
)

// Settings defines the application configuration
type Settings struct {
	InfluxURL      string
	InfluxUser     string
	InfluxPassword string
	InfluxDatabase string
	MQTTHost       string
	MQTTPort       string
	MQTTUser       string
	MQTTPassword   string
	ClientID       string
	UseTLS         bool
	TLSCA          string
	TLSCert        string
	TLSKey         string
	Topics         []string
	TLSCAPem       []byte
	TLSCertPem     []byte
	TLSKeyPem      []byte
}

var settings *Settings

// Read the application configuration
func Read() (*Settings, error) {
	settings = &Settings{}

	// Influxdb settings
	settings.InfluxURL = envVarOrDefault(envInfluxdbURL, defaultInfluxdbURL)
	settings.InfluxDatabase = envVarOrDefault(envInfluxdbDatabase, "")
	settings.InfluxUser = envVarOrDefault(envInfluxdbUser, "")
	settings.InfluxPassword = envVarOrDefault(envInfluxdbPassword, "")

	// MQTT settings
	settings.MQTTHost = envVarOrDefault(envMQTTHost, defaultMQTTHost)
	settings.MQTTPort = envVarOrDefault(envMQTTPort, defaultMQTTPort)
	settings.MQTTUser = envVarOrDefault(envMQTTUser, "")
	settings.MQTTPassword = envVarOrDefault(envMQTTPassword, "")
	settings.ClientID = envVarOrDefault("HOSTNAME", defaultClientID)

	// Topics
	topics := envVarOrDefault(envMQTTTopics, "")
	if len(topics) == 0 {
		log.Fatalf("The list of topics must be supplied: %v", envMQTTTopics)
	}
	settings.Topics = strings.Split(topics, ",")

	// TLS
	tls := envVarOrDefault(envMQTTUseTLS, "")
	switch strings.ToLower(tls) {
	case "true", "t", "yes", "y":
		settings.UseTLS = true
	default:
		settings.UseTLS = false
	}
	settings.TLSCA = envVarOrDefault(envMQTTTLSCA, "")
	settings.TLSCert = envVarOrDefault(envMQTTTLSCert, "")
	settings.TLSKey = envVarOrDefault(envMQTTTLSKEY, "")

	// Read the certs from the file path
	var err error
	if settings.UseTLS {
		settings.TLSCAPem, settings.TLSCertPem, settings.TLSKeyPem, err = readCerts()
		if err != nil {
			log.Printf("Error reading TLS certificates: %v", err)
		}
	}

	return settings, nil
}

func readCerts() ([]byte, []byte, []byte, error) {
	caCert, err := ioutil.ReadFile(settings.TLSCA)
	if err != nil {
		return nil, nil, nil, err
	}
	certCert, err := ioutil.ReadFile(settings.TLSCert)
	if err != nil {
		return nil, nil, nil, err
	}
	certKey, err := ioutil.ReadFile(settings.TLSKey)
	if err != nil {
		return nil, nil, nil, err
	}
	return caCert, certCert, certKey, nil
}

func envVarOrDefault(envVar, defaultValue string) string {
	if len(os.Getenv(envVar)) > 0 {
		return os.Getenv(envVar)
	}
	return defaultValue
}
