package cronjobs

import (
	"log"

	"github.com/robfig/cron/v3"
)

// StartCronJobs initializes and starts all cron jobs
func StartCronJobs() {
	// Create a new Cron scheduler
	c := cron.New()

	// Add a function to check for new alerts every minute
	// _, err := c.AddFunc("@every 1h", func() {

	// 	connections := connectionManager.ConnManager.GetConnectionsByRole(wsTypes.UserRole)

	// 	for _, connection := range connections {
	// 		CheckAndSendAlerts(*connection)
	// 	}

	// })
	// if err != nil {
	// 	log.Fatalf("Error scheduling cron job: %v", err)
	// }

	_, err := c.AddFunc("@every 1m", func() {
		RequestAllSensorData()
	})

	// _, err = c.AddFunc("@every 1m", func() {
	// 	CollectMlData()
	// })

	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	// Start the Cron scheduler
	c.Start()

	log.Println("Cron Job Scheduler started!")
}
