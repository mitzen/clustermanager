package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const slackEndpoint string = "https://hooks.slack.com/services/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"
const applicationJson string = "application/json"

type SlackMessenger struct {
}

func (s *SlackMessenger) SendMessage(message string) {

	postBody, _ := json.Marshal(map[string]string{
		"text": message,
	})

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(slackEndpoint, applicationJson, responseBody)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	sb := string(body)
	log.Printf(sb)
}
