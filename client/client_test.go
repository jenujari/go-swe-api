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
	defer func() { _ = client.Close() }()

	t.Run("Ping", func(t *testing.T) {
		resp, err := client.Ping(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp.Status)
	})

	t.Run("GetPos", func(t *testing.T) {
		resp, err := client.GetPos(t.Context(), "2026-01-26T00:00:00Z", "Sun")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp.Results, "Sun")
	})

	t.Run("Tithy", func(t *testing.T) {
		resp, err := client.Tithy(t.Context(), "2026-01-14T13:45:30Z")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(27), resp.Tithy)
		assert.Equal(t, "Wednesday", resp.Weekday)
		assert.Equal(t, "Anuradha", resp.Nakshatra)
	})

	t.Run("GetBalas", func(t *testing.T) {
		resp, err := client.GetBalas(t.Context(), "2026-01-14T13:45:30Z")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		
		// Assert map has 9 planets
		assert.Len(t, resp.Results, 9)
		
		sun, ok := resp.Results["Sun"]
		assert.True(t, ok)
		assert.NotNil(t, sun)
		assert.InDelta(t, 100.0, sun.UdayBala, 0.001)
		
		// Verify some of the new fields added to PlanetCord
		assert.Equal(t, "Sun", sun.Cords.Name)
		assert.Equal(t, "ati-sheeghra", sun.Cords.SpeedCategory)
		assert.Equal(t, "left-vedha", sun.Cords.Vedha)
		assert.Equal(t, "Saturn", sun.Cords.SignLord)
		assert.Equal(t, "Enemy", sun.Cords.SignLordship)
		assert.NotEmpty(t, sun.Cords.NavamsaSign)
		// Vargottama is boolean, so we just assert it doesn't panic
		_ = sun.Cords.Vargottama
	})
}

func TestFindConjunction(t *testing.T) {
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
	defer func() { _ = client.Close() }()

	resp, err := client.FindConjunction(
		t.Context(),
		"2026-01-01T00:00:00Z",
		"2027-01-01T00:00:00Z",
		"Mars",
		"Saturn",
		1,
		0.041666,
	)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "2026-04-18T10:57:32Z", resp.Start)
	assert.Equal(t, "2026-04-21T10:57:28Z", resp.End)
}
