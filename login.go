package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func login(email string, password string, tokenChan chan string) {
	httpClient, httpReq := getClient(email, password)

	res, err := httpClient.Do(httpReq)
	if err != nil || res.StatusCode != 200 {
		log.Fatal(err)
	}

	var tokenResponse map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		log.Fatal(err)
	}

	currentToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		log.Fatal("access_token not found in response")
	}

	fDuration, ok := tokenResponse["expires_in"].(float64)
	if !ok {
		log.Fatal("expires_in not found in response")
	}
	duration := int64(fDuration)
	tokenExpiration := int64(0)
	tokenExpiration = time.Now().Unix() + duration - 60
	res.Body.Close()

	for {

		if time.Now().Unix() > tokenExpiration {
			httpClient, httpReq := getClient(email, password)
			res, err := httpClient.Do(httpReq)
			if err != nil {
				log.Fatal(err)
			}

			var tokenResponse map[string]string
			err = json.NewDecoder(res.Body).Decode(&tokenResponse)
			if err != nil {
				log.Fatal(err)
			}

			currentToken = tokenResponse["access_token"]

			duration, err := strconv.ParseInt(tokenResponse["expires_in"], 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			tokenExpiration = time.Now().Unix() + duration - 60
			res.Body.Close()
		}

		tokenChan <- currentToken

		time.Sleep(5 * time.Minute)
	}
}

func getClient(email string, password string) (*http.Client, *http.Request) {
	form := url.Values{}
	form.Add("password", password)
	form.Add("username", strings.Replace(email, "@", "%40", 1))
	form.Add("grant_type", "password")

	encodedForm := form.Encode()
	payloadBuffer := bytes.NewBufferString(encodedForm)

	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("POST", "https://api.prod.europe-west1.manual.graduway.com/token", payloadBuffer)
	if err != nil {
		log.Fatal(err)
	}
	httpReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
