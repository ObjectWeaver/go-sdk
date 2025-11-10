package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/objectweaver/go-sdk/converison"
	pb "github.com/objectweaver/go-sdk/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// StreamingResponse represents a single streaming response with converted data
type StreamingResponse struct {
	Data         map[string]any               `json:"data"`
	UsdCost      float64                      `json:"usdCost"`
	Status       string                       `json:"status"`
	DetailedData map[string]*pb.DetailedField `json:"detailedData"`
}

// GrpcStreamGeneratedObjects sends a request to the gRPC server and streams responses
// The handler function is called for each response received from the stream
func (c *Client) GrpcStreamGeneratedObjects(prompt string, definition *pb.Definition, handler func(*StreamingResponse) error) error {
	// Set up a connection to the gRPC server
	conn, err := grpc.NewClient(c.BaseURL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			// Log error but don't return as we're in defer
		}
	}(conn)

	// Create a new client from the gRPC service
	client := pb.NewJSONSchemaServiceClient(conn)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5) // Longer timeout for streaming
	defer cancel()

	// Set up metadata with the authorization token
	md := metadata.New(map[string]string{"x-api-key": c.Password})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Create the request object
	request := &pb.RequestBody{
		Prompt:     prompt,
		Definition: definition,
	}

	// Call the streaming gRPC method
	stream, err := client.StreamGeneratedObjects(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to call StreamGeneratedObjects: %v", err)
	}

	// Process the stream
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// Stream completed successfully
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %v", err)
		}

		// Convert the structpb data to map
		data, err := converison.ConvertStructpbToMap(response.Data)
		if err != nil {
			return fmt.Errorf("error converting response data: %v", err)
		}

		// Create streaming response
		streamResp := &StreamingResponse{
			Data:         data,
			UsdCost:      response.UsdCost,
			Status:       response.Status,
			DetailedData: response.DetailedData,
		}

		// Call the handler function
		if err := handler(streamResp); err != nil {
			return fmt.Errorf("handler error: %v", err)
		}
	}

	return nil
}
