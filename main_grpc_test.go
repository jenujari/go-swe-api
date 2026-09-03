package main

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/jenujari/go-swe-api/grpc"
	pb "github.com/jenujari/go-swe-api/proto"
	"github.com/stretchr/testify/assert"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis          *bufconn.Listener
	testGRPCOnce sync.Once
)

func initTestGRPC() {
	testGRPCOnce.Do(func() {
		lis = bufconn.Listen(bufSize)
		s := googlegrpc.NewServer()
		pb.RegisterEphServiceServer(s, &grpc.Server{})
		go func() {
			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		}()
	})
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func grpcClient(t *testing.T) (context.Context, pb.EphServiceClient, func()) {
	t.Helper()
	initTestGRPC()
	ctx := context.Background()
	conn, err := googlegrpc.NewClient(
		"passthrough:///bufnet",
		googlegrpc.WithContextDialer(bufDialer),
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return ctx, pb.NewEphServiceClient(conn), func() { _ = conn.Close() }
}

func TestGRPC_Ping(t *testing.T) {
	ctx, client, closeConn := grpcClient(t)
	defer closeConn()

	resp, err := client.Ping(ctx, &pb.PingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestGRPC_GetPos(t *testing.T) {
	ctx, client, closeConn := grpcClient(t)
	defer closeConn()

	resp, err := client.GetPos(ctx, &pb.PosRequest{
		Time:       "2026-01-26T00:00:00Z",
		PlanetName: "Sun",
	})
	assert.NoError(t, err)
	sun, ok := resp.Results["Sun"]
	assert.True(t, ok)
	assert.NotNil(t, sun)
	assert.NotEmpty(t, sun.VedhaTarget)
	assert.NotNil(t, sun.LongitudeDms)
}

func TestGRPC_GetBalas(t *testing.T) {
	ctx, client, closeConn := grpcClient(t)
	defer closeConn()

	resp, err := client.GetBalas(ctx, &pb.BalasRequest{
		Timestamp: "2026-01-14T13:45:30Z",
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Results, 9)
	sun, ok := resp.Results["Sun"]
	assert.True(t, ok)
	assert.NotNil(t, sun.Cords)
	assert.NotEmpty(t, sun.Cords.VedhaTarget)
}

func TestGRPC_Tithy(t *testing.T) {
	ctx, client, closeConn := grpcClient(t)
	defer closeConn()

	resp, err := client.Tithy(ctx, &pb.TithyRequest{
		Timestamp: "2026-01-14T13:45:30Z",
	})
	assert.NoError(t, err)
	assert.NotZero(t, resp.Tithy)
	assert.NotEmpty(t, resp.Nakshatra)
	assert.NotEmpty(t, resp.Weekday)
}

func TestGRPC_FindConjunction(t *testing.T) {
	ctx, client, closeConn := grpcClient(t)
	defer closeConn()

	resp, err := client.FindConjunction(ctx, &pb.ConjunctionRequest{
		Start:   "2026-01-01T00:00:00Z",
		End:     "2026-03-02T00:00:00Z",
		Planet1: "Sun",
		Planet2: "Mercury",
		Orb:     1,
		Step:    1.0 / 24.0,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Start)
	assert.NotEmpty(t, resp.End)
}
