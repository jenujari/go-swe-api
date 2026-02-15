package main

import (
	"context"
	"net"
	"testing"

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
	pb.RegisterSWEServiceServer(s, &grpc.Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGRPC_Ping(t *testing.T) {
	initTestGRPC()
	ctx := context.Background()
	conn, err := googlegrpc.DialContext(ctx, "bufnet", googlegrpc.WithContextDialer(bufDialer), googlegrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSWEServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
	assert.Equal(t, "ok", resp.Status)
}

func TestGRPC_GetPos(t *testing.T) {
	// initTestGRPC() // Already initialized if running all tests, but better safe.
	ctx := context.Background()
	conn, err := googlegrpc.DialContext(ctx, "bufnet", googlegrpc.WithContextDialer(bufDialer), googlegrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSWEServiceClient(conn)
	resp, err := client.GetPos(ctx, &pb.PosRequest{
		Time:       "2026-01-26T00:00:00Z",
		PlanetName: "Sun",
	})
	if err != nil {
		t.Fatalf("GetPos failed: %v", err)
	}

	sun, ok := resp.Results["Sun"]
	assert.True(t, ok)
	assert.NotNil(t, sun)
	assert.InDelta(t, 281.808299, sun.Longitude, 0.001)
}
