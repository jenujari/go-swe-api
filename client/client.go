package client

import (
	"context"
	"fmt"

	pb "github.com/jenujari/go-swe-api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EphServiceClient is a wrapper around the generated gRPC client.
type EphServiceClient struct {
	conn   *grpc.ClientConn
	client pb.EphServiceClient
}

// NewSWEClient creates a new SWEClient and establishes a connection to the gRPC server.
// The address should be in the format "host:port".
func NewEphServiceClient(address string) (*EphServiceClient, error) {
	// Using NewClient instead of Dial (which is deprecated)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server at %s: %w", address, err)
	}

	return &EphServiceClient{
		conn:   conn,
		client: pb.NewEphServiceClient(conn),
	}, nil
}

// Close closes the underlying gRPC connection.
func (c *EphServiceClient) Close() error {
	return c.conn.Close()
}

// Ping checks the health of the gRPC server.
func (c *EphServiceClient) Ping(ctx context.Context) (*pb.PingResponse, error) {
	return c.client.Ping(ctx, &pb.PingRequest{})
}

// GetPos retrieves the position of a planet at a specific time.
// If planetName is empty, it returns positions for all available planets.
// timeStr should be in RFC3339 format (e.g., "2026-01-26T00:00:00Z").
func (c *EphServiceClient) GetPos(ctx context.Context, timeStr string, planetName string) (*pb.PosResponse, error) {
	return c.client.GetPos(ctx, &pb.PosRequest{
		Time:       timeStr,
		PlanetName: planetName,
	})
}

// FindConjunction searches for conjunctions between two planets within a time range.
func (c *EphServiceClient) FindConjunction(ctx context.Context, start, end, planet1, planet2 string, orb int32, step float64) (*pb.ConjunctionResponse, error) {
	return c.client.FindConjunction(ctx, &pb.ConjunctionRequest{
		Start:   start,
		End:     end,
		Planet1: planet1,
		Planet2: planet2,
		Orb:     orb,
		Step:    step,
	})
}
