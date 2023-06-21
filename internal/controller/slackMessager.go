package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const applicationJson string = "application/json"

type SlackMessenger struct {
	endpoint string
}

func NewSlackMessenger(endpoint string) *SlackMessenger {
	return &SlackMessenger{endpoint: endpoint}
}

func (sm *SlackMessenger) SendMessage(message string) int16 {

	postBody, _ := json.Marshal(map[string]string{
		"text": message,
	})

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(sm.endpoint, applicationJson, responseBody)

	if err != nil {
		log.Fatalf("An error occured %v", err)
		return 1
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("An error occured %v deserializing the response message", err)
		return 2
	}

	sb := string(body)
	log.Printf("%s", sb)

	return 0
}
