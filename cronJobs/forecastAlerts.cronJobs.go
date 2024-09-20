package cronjobs

import (
	"PlantCare/websocket/connectionManager"
	"log"
	"github.com/robfig/cron/v3"
)

// StartCronJobs initializes and starts all cron jobs
func StartCronJobs() {
	// Create a new Cron scheduler
	c := cron.New()

	// Add a function to check for new alerts every minute
	_, err := c.AddFunc("@every 1m", func() {


		connection, ok := connectionManager.ConnManager.GetConnection("user_2jod4hRuJ9nqUIzftpaTTNWLVxv")
        if !ok {
            return
        }

        if connection == nil{
            return
        }

        CheckAndSendAlerts(*connection)

	})

	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	// Start the Cron scheduler
	c.Start()

	log.Println("Cron Job Scheduler started!")
}
