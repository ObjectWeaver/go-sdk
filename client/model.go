package client

import "github.com/firechimp-org/go-sdk/jsonSchema"

type RequestBody struct {
	Prompt     string                 `json:"prompt"`
	Definition *jsonSchema.Definition `json:"definition"`
}

// Create a response struct
type Response struct {
	Data    map[string]any `json:"data"` //this data can then be marshalled into the apprioate object type.
	UsdCost float64        `json:"usdCost"`
}
