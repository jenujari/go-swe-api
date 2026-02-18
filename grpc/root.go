package grpc

import (
	"fmt"
	"net"

	"github.com/jenujari/go-swe-api/config"
	pb "github.com/jenujari/go-swe-api/proto"
	rtc "github.com/jenujari/runtime-context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcServer *grpc.Server
)

func init() {
	grpcServer = grpc.NewServer()
	pb.RegisterEphServiceServer(grpcServer, &Server{})
	reflection.Register(grpcServer)
}

func RunGRPCServer() {
	pc := rtc.GetMainProcess()
	cfg := config.GetConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.Port))
	if err != nil {
		pc.FatalErrorChan <- fmt.Errorf("failed to listen: %v", err)
		return
	}

	config.GetLogger().Println("gRPC server initialization complete.")

	go func(cmdx *rtc.ProcessContext) {
		config.GetLogger().Println("gRPC Server is running at", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			cmdx.FatalErrorChan <- fmt.Errorf("gRPC Serve(): %v", err)
		}
	}(pc)

	<-pc.CTX.Done()
	config.GetLogger().Println("shutting down gRPC server...")
	grpcServer.GracefulStop()
	config.GetLogger().Println("gRPC server shutdown complete...")
}
