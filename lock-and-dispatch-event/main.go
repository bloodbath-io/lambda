package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

const (
	callbackEndpoint = "https://api.bloodbath.io/internal/callback"
)

type Payload struct {
	Id       string `json:"id"`
	Body     string `json:"body"`
	Endpoint string `json:"endpoint"`
	Headers  string `json:"headers"`
	Method   string `json:"method"`
}

type Response struct {
	Id     string
	Type   string
	Status int
	Body   string
	Reason string
}

func handleRequest(context context.Context, payload Payload) error {
	response, err := sendRequest(context, payload)
	if err != nil {
		throwError(err, payload)
	}

	fmt.Printf("Response body: %s\r\n", string(response.Body))
	fmt.Printf("Response status: %s\r\n", fmt.Sprint(response.Status))

	err = sendCallback(response)
	if err != nil {
		throwError(err, payload)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}

func throwError(err error, payload Payload) {
	sendCallback(Response{Type: "error", Id: payload.Id, Reason: err.Error()})
	log.Fatal(err)
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
	if err != nil {
		return Response{}, err
	}

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

	return Response{Type: "ok", Id: id, Status: resp.StatusCode, Body: string(responseBody)}, nil
}

func sendCallback(response Response) error {
	body := &Response{
		Type:   response.Type,
		Id:     response.Id,
		Body:   response.Body,
		Status: response.Status,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	request, err := http.NewRequest("POST", callbackEndpoint, payloadBuf)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	endResponse, err := client.Do(request)
	if err != nil {
		return err
	}

	defer endResponse.Body.Close()

	fmt.Println("response Status:", endResponse.StatusCode)
	fmt.Println("response Headers:", endResponse.Header)
	endBody, err := ioutil.ReadAll(endResponse.Body)
	if err != nil {
		return err
	}

	fmt.Println("response Body:", string(endBody))

	return nil
}
