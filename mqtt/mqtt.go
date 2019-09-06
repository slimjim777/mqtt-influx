package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/slimjim777/mqtt-influx/config"
	"github.com/slimjim777/mqtt-influx/datastore"
	"log"
)

// Constants for connecting to the MQTT broker
const (
	quiesce        = 250
	QOSAtMostOnce  = byte(0)
	QOSAtLeastOnce = byte(1)
	//QOSExactlyOnce = byte(2)
)

// Connection for MQTT protocol
type Connection struct {
	client   MQTT.Client
	clientID string
	settings *config.Settings
}

var conn *Connection
var client MQTT.Client

// GetConnection fetches or creates an MQTT connection
func GetConnection(settings *config.Settings) (*Connection, error) {
	if conn == nil {
		// Create the client
		client, err := newClient(settings)
		if err != nil {
			return nil, err
		}

		// Create a new connection
		conn = &Connection{
			client:   client,
			clientID: settings.ClientID,
			settings: settings,
		}
	}

	// Check that we have a live connection
	if conn.client.IsConnectionOpen() {
		return conn, nil
	}

	// Connect to the MQTT broker
	if token := conn.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	// Subscribe to the actions topic
	err := conn.SubscribeToActions()
	return conn, err
}

// Close the connection
func (conn *Connection) Close() {
	conn.Close()
}

func (conn *Connection) SubscribeHandler(client MQTT.Client, msg MQTT.Message) {
	log.Println("Write message:", string(msg.Payload()))
	if err := datastore.Write(conn.settings, string(msg.Payload())); err != nil {
		log.Println("Error storing message:", err)
	}
}

// SubscribeToActions subscribes to the topics
func (conn *Connection) SubscribeToActions() error {
	log.Println("Subscribe to actions...")
	for _, t := range conn.settings.Topics {
		log.Println("...action:", t)
		token := client.Subscribe(t, QOSAtLeastOnce, conn.SubscribeHandler)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Error subscribing to topic `%s`: %v", t, token.Error())
			return fmt.Errorf("error subscribing to topic `%s`: %v", t, token.Error())
		}
	}

	return nil
}

// newClient creates a new MQTT client
func newClient(settings *config.Settings) (MQTT.Client, error) {
	var url string
	// Return the active client, if we have one
	if client != nil {
		return client, nil
	}

	// Generate a new MQTT client
	if settings.UseTLS {
		url = fmt.Sprintf("ssl://%s:%s", settings.MQTTHost, settings.MQTTPort)
	} else {
		url = fmt.Sprintf("tcp://%s:%s", settings.MQTTHost, settings.MQTTPort)
	}
	log.Println("Connect to the MQTT broker", url)

	// Set up the MQTT client options
	opts := MQTT.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID(settings.ClientID)

	// Set the TLS certs, if needed
	if settings.UseTLS {
		// Generate the TLS config from the enrollment credentials
		tlsConfig, err := newTLSConfig(settings)
		if err != nil {
			return nil, err
		}
		opts.SetTLSConfig(tlsConfig)
	}

	// Set the username and password, if needed
	if len(settings.MQTTUser) > 0 {
		opts.SetUsername(settings.MQTTUser)
		opts.SetPassword(settings.MQTTPassword)
	}

	// Client to reconnect on disconnect
	opts.AutoReconnect = true

	client = MQTT.NewClient(opts)
	return client, nil
}

// newTLSConfig sets up the certificates from the enrollment record
func newTLSConfig(settings *config.Settings) (*tls.Config, error) {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(settings.TLSCAPem)

	// Import client certificate/key pair
	cert, err := tls.X509KeyPair(settings.TLSCertPem, settings.TLSKeyPem)
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired TLS properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certPool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}
