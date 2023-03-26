package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var adbClient = *AdbConnect()

func HandleRequest() {

	API_URL := "https://threatfox-api.abuse.ch/api/v1/"

	data := map[string]interface{}{"query": "get_iocs", "days": 1}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	// Send the request
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Parse the response
	var threatFoxResponse ThreatFoxResponse
	err = json.Unmarshal(body, &threatFoxResponse)
	if err != nil {
		log.Fatalf("Error parsing response JSON: %v", err)
	}

	vertexCollection, edgeCollection := BuildCollections(threatFoxResponse)

	for collectionName, collection := range vertexCollection {
		EnsureDocuments(context.TODO(), collectionName, collection)
	}

	for collectionName, collection := range edgeCollection {
		EnsureEdgeDocuments(context.TODO(), collectionName, collection)
	}

	log.Print("Completed ThreatFox data sync.")
}

func main() {
	lambda.Start(HandleRequest)
}
