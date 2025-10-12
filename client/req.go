package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/firechimp-org/go-sdk/jsonSchema"
	"log"
	"net/http"
)

// ExecuteRequest executes an HTTP request based on the RequestFormat
func ExecuteRequest(currentGen map[string]any, d *jsonSchema.Definition) (*http.Response, error) {
	// Merge currentGen into the existing body
	if d.Req.Body == nil {
		d.Req.Body = make(map[string]interface{})
	}
	for key, value := range currentGen {
		d.Req.Body[key] = value
	}

	// Marshal the body to JSON
	body, err := json.Marshal(d.Req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(string(d.Req.Method), d.Req.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range d.Req.Headers {
		req.Header.Set(key, value)
	}
	// Set authorization
	req.Header.Set("Authorization", d.Req.Authorization)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// the below occurs within the internal workings (i think)
func SendRequest(def *jsonSchema.Definition, currentGen map[string]any) *Res {
	//the request that is sent out needs to send out a current map of the generated object.

	request, err := ExecuteRequest(currentGen, def)
	if err != nil {
		log.Println("failed to execute request", err)
		return nil
	}

	value, err := extractValue(request)
	if err != nil {
		log.Println("failed to extract value", err)
		return nil
	}

	return value
}
