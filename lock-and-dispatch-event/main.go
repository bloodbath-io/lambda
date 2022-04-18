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

type Headers struct {
	ContentType string `json:"content-Type"`
}

func handleRequest(context context.Context, payload Payload) error {
	printablePayload, _ := json.Marshal(payload)
	printableContext, _ := json.Marshal(context)
	fmt.Printf("Received context %s\r\n", string(printableContext))
	fmt.Printf("Received payload %s\r\n", string(printablePayload))

	id := payload.Id
	body, _ := json.Marshal(payload.Body)
	endpoint := payload.Endpoint

	// 	var formattedHeaders map[string]string

	var headers Headers
	json.Unmarshal([]byte(payload.Headers), &headers)

	// headers, _ := json.Marshal(payload.Headers)
	// _ = json.Unmarshal(headers, &formattedHeaders)

	method := strings.ToUpper(payload.Method)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))

	fmt.Printf("Content-Type: %s\r\n", headers.ContentType)
	req.Header.Set("Content-Type", headers.ContentType)

	// we need to transmit the host
	// as a header to avoid 400 Bad Request
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Host: %s\r\n", parsedEndpoint.Host)
	req.Host = parsedEndpoint.Host

	// other headers
	req.Header.Set("Cache-Control", "no-cache")

	// fmt.Printf("This is the headers we will send: %s", headers)

	// for key, value := range headers {
	// 	fmt.Printf("Parsing the header `%s` : `%s`\r\n", key, value)
	// 	req.Header.Set(key, value)
	// }

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // appending to existing query args
	// query := req.URL.Query()
	// // q.Add("foo", "bar")

	// // assign encoded query string to http request
	// req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
		return err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: deal with response and send it back to bloodbath
	fmt.Printf("It went well for ID: %s", id)
	fmt.Printf("Status received is: %s", resp.Status)
	fmt.Printf("Body received is: %s", string(responseBody))

	// TESTING OUT
	// hey, err := http.NewRequest("GET", "https://dummy.bloodbath.io/lets-try", nil)
	// hey.Host = "dummy.bloodbath.io"
	// clientx := &http.Client{}
	// yo, err := clientx.Do(hey)

	// // yo, err := http.Get("https://dummy.bloodbath.io/lets-try")
  // // if err != nil {
  // //     log.Fatal(err)
  // // }

  // defer yo.Body.Close()

  // b, err := ioutil.ReadAll(yo.Body)
  // if err != nil {
  //     log.Fatal(err)
  // }

	// fmt.Printf("TESTING AFTER THAT \r\n")
  // fmt.Printf(string(b))
	// END OF TESTING OUT

	return nil
}

func main() {
	lambda.Start(handleRequest)
}

// cannot unmarshal string into Go value of type lambdacontext.ClientContext
