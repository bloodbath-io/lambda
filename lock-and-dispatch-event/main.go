package main

import (
	"bytes"
	"strings"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Payload struct {
	Id       string `json:"id"`
	Body     string `json:"body"`
	Endpoint string `json:"endpoint"`
	Headers  string `json:"headers"`
	Method   string `json:"method"`
}

type Response struct {
	Status string
	Body string
}

func handleRequest(context context.Context, payload Payload) error {
	response, err := sendRequest(context, payload)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response body: %s\r\n", string(response.Body))
	fmt.Printf("Response status: %s\r\n", string(response.Status))

	return nil
}

func main() {
	lambda.Start(handleRequest)
}

func sendRequest(context context.Context, payload Payload) (Response, error) {
	printablePayload, _ := json.Marshal(payload)
	printableContext, _ := json.Marshal(context)
	fmt.Printf("Received context %s\r\n", string(printableContext))
	fmt.Printf("Received payload %s\r\n", string(printablePayload))

	id := payload.Id
	body, _ := json.Marshal(payload.Body)
	endpoint := payload.Endpoint

	var headers map[string]string
	json.Unmarshal([]byte(payload.Headers), &headers)

	method := strings.ToUpper(payload.Method)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))

	// default headers are added
	req.Header.Set("Cache-Control", "no-cache")

	// we need to transmit the host
	// as a header to avoid 400 Bad Request
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return Response{}, err
	}

	fmt.Printf("Host: %s\r\n", parsedEndpoint.Host)
	req.Host = parsedEndpoint.Host

	// we parse all the headers sent
	// from bloodbath and set them up
	for key, value := range headers {
		fmt.Printf("Parsing the header `%s` : `%s`\r\n", key, value)
		req.Header.Set(key, value)
	}

	if err != nil {
		return Response{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
		return Response{}, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{Status: resp.Status, Body: string(responseBody)}, nil
}