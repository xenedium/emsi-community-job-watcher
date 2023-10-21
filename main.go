package main

import (
	"log"
	"os"
)

func main() {
	email := os.Getenv("EC_EMAIL")
	password := os.Getenv("EC_PASSWORD")
	funcURL := os.Getenv("EC_FUNC_URL")
	discordWebhook := os.Getenv("EC_DISCORD_WEBHOOK")

	if email == "" || password == "" {
		log.Fatal("EC_EMAIL and EC_PASSWORD must be set")
	}

	if funcURL == "" && discordWebhook == "" {
		log.Fatal("EC_FUNC_URL or EC_DISCORD_WEBHOOK must be set")
	}

	tokenChan := make(chan string)
	go login(email, password, tokenChan) // first goroutine to keep the token updated if it expires

	jobsChan := make(chan Job)
	go extractJobs(tokenChan, jobsChan) // second goroutine to extract jobs from the API

	if funcURL != "" {
		go sendJobsToFunction(funcURL, jobsChan) // third goroutine to send jobs to the function
	}
	if discordWebhook != "" {
		go sendJobsToDiscord(discordWebhook, jobsChan) // fourth goroutine to send jobs to Discord
	}

	// main goroutine to keep the program running
	select {}

}
