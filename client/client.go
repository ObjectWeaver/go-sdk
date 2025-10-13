package client

import (
	"net/http"

	"github.com/objectweaver/go-sdk/jsonSchema"
)

// Client responsible for holding base configuration and dependencies
type Client struct {
	Password          string
	BaseURL           string
	HttpClient        HttpClient
	RequestSender     RequestSender
	ResponseProcessor ResponseProcessor
}

// HttpClient interface to abstract HTTP operations
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RequestSender interface abstracts request sending behavior
type RequestSender interface {
	SendRequestBody(url, token string, requestBody *RequestBody) (*http.Response, error)
}

// NewDefaultClient initializes a new Client instance with default implementations
func NewDefaultClient(password, url string, client *http.Client) *Client {
	return &Client{
		Password:          password,
		BaseURL:           url,
		HttpClient:        client,
		RequestSender:     NewDefaultRequestSender(client),
		ResponseProcessor: NewResponseProcessor(),
	}
}

// NewGZipClient initializes a new Client instance with GZip compression for requests
func NewGZipClient(password, url string, client *http.Client) *Client {
	return &Client{
		Password:          password,
		BaseURL:           url,
		HttpClient:        client,
		RequestSender:     NewGZipRequestSender(client),
		ResponseProcessor: NewResponseProcessor(),
	}
}

// SendRequest sends the prompt and definition, and returns the parsed response
func (c *Client) SendRequest(prompt string, definition *jsonSchema.Definition) (*Response, error) {
	requestBody := &RequestBody{
		Prompt:     prompt,
		Definition: definition,
	}

	// Use the RequestSender to send the request
	resp, err := c.RequestSender.SendRequestBody(c.BaseURL, c.Password, requestBody)
	if err != nil {
		return nil, err
	}

	// Process the response
	return c.ResponseProcessor.ProcessResponse(resp)
}
