package client

import (
	"context"
	"fmt"
	"github.com/firechimp-org/go-sdk/converison"
	pb "github.com/firechimp-org/go-sdk/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// SendRequestToServer sends a request to the gRPC server with authorization headers
func (c *Client) GrpcGenerateObject(prompt string, definition *pb.Definition) (*Response, error) {
	// Set up a connection to the gRPC server
	conn, err := grpc.NewClient(c.BaseURL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {

		}
	}(conn)

	// Create a new client from the gRPC service
	client := pb.NewJSONSchemaServiceClient(conn)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Set up metadata with the authorization token
	md := metadata.New(map[string]string{"x-api-key": c.Password})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Create the request object
	request := &pb.RequestBody{
		Prompt:     prompt,
		Definition: definition,
	}

	// Call the gRPC method on the client
	response, err := client.GenerateObject(ctx, request) // Replace 'GenerateObject' with the actual RPC method name
	if err != nil {
		return nil, fmt.Errorf("failed to call GenerateObject: %v", err)
	}

	data, err := converison.ConvertStructpbToMap(response.Data)

	res := &Response{
		Data:    data,
		UsdCost: response.UsdCost,
	}

	return res, nil
}
