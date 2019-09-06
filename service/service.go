package control

import (
	"github.com/slimjim777/mqtt-influx/config"
	"log"
	"time"
)

const sleepMinutes = 2

// Run loop for the service
func Run(settings *config.Settings) error {
	for {
		log.Println("Check")

		// Wait before repeat
		log.Printf("Wait for %d minutes before restarting", sleepMinutes)
		time.Sleep(sleepMinutes * time.Minute)
	}

	return nil
}
