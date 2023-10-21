package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func sendJobsToFunction(funcURL string, jobs chan Job) {
	for job := range jobs {
		jsonJob, err := json.Marshal(job)
		if err != nil {
			continue
		}
		http.Post(funcURL, "application/json", bytes.NewBuffer(jsonJob))
	}
}

func sendJobsToDiscord(discordWebhook string, jobs chan Job) {
	registredJobs := make([]string, 0)
	for job := range jobs {
		if !contains(registredJobs, job.ID) {
			registredJobs = append(registredJobs, job.ID)

			discordEmbed := make(map[string]interface{})
			discordEmbed["title"] = job.Title
			discordEmbed["url"] = job.Link
			discordEmbed["description"] = job.Description
			discordEmbed["color"] = 16777215
			discordEmbed["timestamp"] = job.Posted
			discordEmbed["footer"] = map[string]string{"text": job.Company}
			discordEmbed["thumbnail"] = map[string]string{"url": job.CompanyLogo}
			discordEmbed["author"] = map[string]string{"name": job.Location}
			discordPayload := make(map[string]interface{})
			discordPayload["embeds"] = []map[string]interface{}{discordEmbed}
			jsonPayload, err := json.Marshal(discordPayload)
			if err != nil {
				continue
			}
			http.Post(discordWebhook, "application/json", bytes.NewBuffer(jsonPayload))
		}
	}
}

func contains(registredJobs []string, s string) bool {
	for _, a := range registredJobs {
		if a == s {
			return true
		}
	}
	return false
}
