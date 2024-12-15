package main

import (
	"log"
	"resume_website/app"
	"time"

	"github.com/robfig/cron/v3"
)

func scheduleDailyRestart() {
	c := cron.New()
	_, err := c.AddFunc("0 1 * * *", func() {
		log.Println("Restarting server...")
		app.StopServer()
		time.Sleep(10 * time.Second)
		go app.StartServer()
	})
	if err != nil {
		log.Fatalf("Error scheduling daily restart: %v", err)
	}
	c.Start()
}

func main() {
	app.StartServer()
}
