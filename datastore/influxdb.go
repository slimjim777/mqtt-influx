package datastore

import (
	"github.com/slimjim777/mqtt-influx/config"
	"net/http"
	"net/url"
	"strings"
)

// Write outputs a point
func Write(settings *config.Settings, data string) error {
	// Set up the URL
	u, err := url.Parse(settings.InfluxURL)
	if err != nil {
		return err
	}
	u.Path = "write"
	q := u.Query()
	q.Set("db", settings.InfluxDatabase)

	// Authentication
	if len(settings.InfluxUser) > 0 {

		q.Set("username", settings.InfluxUser)
		q.Set("password", settings.InfluxPassword)
	}
	u.RawQuery = q.Encode()

	_, err = http.Post(u.String(), "", strings.NewReader(data))
	return err
}
