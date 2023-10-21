package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func extractJobs(token chan string, jobs chan Job) {

	for {
		httpClient, httpReq := getHttpClient()
		receiveToken := <-token
		log.Printf("Received token, starting to extract jobs")
		httpReq.Header.Add("Authorization", "bearer "+receiveToken)

		res, err := httpClient.Do(httpReq)
		if err != nil || res.StatusCode != 200 {
			log.Fatal(err)
		}

		var jobsResponse map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&jobsResponse)
		if err != nil {
			log.Fatal(err)
		}

		jobsResponseContent, ok := jobsResponse["content"].(map[string]interface{})
		if !ok {
			log.Fatal("content not found in response")
		}
		jobsResponseData1, ok := jobsResponseContent["data"].(map[string]interface{})
		if !ok {
			log.Fatal("data not found in response")
		}
		jobsResponseData2, ok := jobsResponseData1["data"].([]interface{})
		if !ok {
			log.Fatal("data not found in response")
		}

		for _, job := range jobsResponseData2 {
			newJob := Job{}
			jobMap, ok := job.(map[string]interface{})
			if !ok {
				continue
			}

			id, ok := jobMap["id"].(float64)
			if !ok {
				continue
			}
			newJob.ID = strconv.FormatInt(int64(id), 10)
			newJob.Title, _ = jobMap["title"].(string)
			company, ok := (jobMap["company"].(map[string]interface{}))
			if ok {
				newJob.Company, _ = company["name"].(string)
				newJob.CompanyLogo, _ = company["logoUrl"].(string)
			}
			location, ok := (jobMap["location"].(map[string]interface{}))
			if ok {
				newJob.Location, _ = location["name"].(string)
			}
			newJob.Posted, _ = jobMap["postDate"].(string)
			newJob.Link = "https://emsicommunity.com/jobs/" + newJob.ID
			newJob.Description, _ = jobMap["description"].(string)

			jobs <- newJob
		}
		res.Body.Close()
	}
}

func getHttpClient() (*http.Client, *http.Request) {
	jsonString := "{\"offsetRowsToSkip\":0,\"maxRowsToFetch\":6,\"companyName\":\"\",\"jobTitle\":\"\",\"jobTypes\":[7,8],\"jobFunctions\":[],\"industries\":[],\"location\":null,\"locationFields\":{},\"freeText\":\"\",\"jobArchiveType\":0}"
	payloadBuffer := bytes.NewBufferString(jsonString)

	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("POST", "https://api.prod.europe-west1.manual.graduway.com/Jobs/search/", payloadBuffer)
	if err != nil {
		log.Fatal(err)
	}
	httpReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Access-Control-Allow-Origin", "*")
	httpReq.Header.Add("Origin", "https://emsicommunity.com")
	httpReq.Header.Add("Referer", "https://emsicommunity.com/")
	httpReq.Header.Add("HorizontalId", "50152")
	httpReq.Header.Add("HorizontalName", "emsicommunity")
	httpReq.Header.Add("SharedLanguageId", "1154")
	httpReq.Header.Add("onlyAdminAllowed", "false")
	httpReq.Header.Add("Sec-Fetch-Dest", "empty")
	httpReq.Header.Add("Sec-Fetch-Mode", "cors")
	httpReq.Header.Add("Sec-Fetch-Site", "cross-site")
	httpReq.Header.Add("Sec-GPC", "1")

	return httpClient, httpReq
}
