package client

import (
	"context"
	"net"
	"testing"

	"github.com/jenujari/go-swe-api/config"
	"github.com/jenujari/go-swe-api/grpc"
	pb "github.com/jenujari/go-swe-api/proto"
	"github.com/stretchr/testify/assert"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func initTestGRPC() {
	lis = bufconn.Listen(bufSize)
	s := googlegrpc.NewServer()
	pb.RegisterEphServiceServer(s, &grpc.Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestEphServiceClient(t *testing.T) {

	// Optionally set a mock config
	config.SetConfig(&config.Config{
		App: struct {
			Name  string `mapstructure:"name"`
			Port  int    `mapstructure:"port"`
			Debug bool   `mapstructure:"debug"`
		}{
			Name:  "test-app",
			Port:  5678,
			Debug: true,
		},
	})

	initTestGRPC()
	ctx := context.Background()

	// Create client using the new EphServiceClient wrapper
	// We use DialContext with bufDialer to test against the in-memory server
	opts := []googlegrpc.DialOption{
		googlegrpc.WithContextDialer(bufDialer),
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := googlegrpc.NewClient("passthrough:///bufnet", opts...)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := &EphServiceClient{
		conn:   conn,
		client: pb.NewEphServiceClient(conn),
	}
	defer client.Close()

	t.Run("Ping", func(t *testing.T) {
		resp, err := client.Ping(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp.Status)
	})

	t.Run("GetPos", func(t *testing.T) {
		resp, err := client.GetPos(ctx, "2026-01-26T00:00:00Z", "Sun")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp.Results, "Sun")
	})
}
