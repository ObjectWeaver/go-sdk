# Object Generation Golang

<p align="center">
  <img src="https://i.postimg.cc/K8TjGy01/3-BD4-B9-DD-B0-B4-4033-B561-50-B5-CB6-A37-B0.jpg" alt="Alt text" width="200"/>
</p>

This module provides a client implementation for sending JSON definitions via HTTP POST requests. It is designed to be simple and easy to integrate into your existing Go projects.

## Installation

```go get github.com/henrylamb/object-generation-golang```

To use this library you will need to import it in the below format:

```go
import (
"github.com/henrylamb/object-generation-golang/jsonSchema"
)

```

## Guide: Using Go Client to Send JSON Definitions

This guide demonstrates how to create a Go client that sends JSON definitions using HTTP POST requests.

### Step 1: Define the Client Struct

Start by defining a `Client` struct that will manage the API connection.

```go
type Client struct {
	Password  string // APIKey is the authentication token for the API.
	BaseURL string // BaseURL is the base endpoint for API requests.
}
```

# Guide: Using Go Client to Send JSON Definitions

This guide demonstrates how to create a Go client that sends JSON definitions using HTTP POST requests.

### Step 1: Define the Client Struct

Start by defining a `Client` struct that will manage the API connection.

```go
type Client struct {
    Password   string
    BaseURL    string
    HttpClient HttpClient
}
```

### Step 2: Initialize a New Client

Create a constructor function `NewClient` to initialize a new client instance with the API key and base URL.

```go
func NewDefaultClient(password, url string) *Client {
	return NewClient(password, url, &http.Client{})
}

// NewClient initializes a new Client instance
func NewClient(password, url string, httpClient HttpClient) *Client {
	return &Client{
		Password:   password,
		BaseURL:    url,
		HttpClient: httpClient,
	}
}
```

### Step 3: Define the SendRequest Method

Implement a method `SendRequest` on the `Client` struct to send a POST request with a JSON-encoded definition.

```go
// SendRequest sends the prompt and definition, and returns the parsed response
func (c *Client) SendRequest(prompt string, definition *jsonSchema.Definition) (*Response, error) {
	requestSender := NewRequestSender(c)
	responseProcessor := NewResponseProcessor()

	requestBody := &RequestBody{
		Prompt:     prompt,
		Definition: definition,
	}

	resp, err := requestSender.SendRequestBody(requestBody)
	if err != nil {
		return nil, err
	}

	return responseProcessor.ProcessResponse(resp)
}
```

### Step 4: Example Usage

Demonstrate how to use the `Client` to send a definition using a sample `Definition` struct.

```go
// Example usage
func ExampleUsage() {
	// Initialize a new client with your API key
	url   := "your-container-url"
	password := "your-password"
	c   := client.NewDefaultClient(password, url)

	// Define a sample definition
	definition := &Definition{
		Type:        "Object",
		Instruction: "Sample instruction for the definition.",
		Properties: map[string]Definition{
			"property1": {Type: "String", Instruction: "Description of property1"},
			"property2": {Type: "Number", Instruction: "Description of property2"},
		},
	}

	// Send the request
	resp, err := c.SendRequest(definition)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Process response as needed
	fmt.Printf("Response Status: %s\n", resp.Status)
	// Additional processing of response body, headers, etc.
}
```

### Conclusion

This guide provides a structured approach to creating a Go client for sending JSON definitions via HTTP POST requests. Ensure to adapt the `Definition` struct and example usage to fit your specific API requirements and data structures.
