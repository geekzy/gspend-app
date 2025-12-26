package grpc

import (
	"fmt"
	"net"

	authv1 "github.com/geekzy/gspend-app/apps/auth-service/internal/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	port   int
}

func NewGRPCServer(authGRPCService *AuthGRPCService, port int) *GRPCServer {
	server := grpc.NewServer()
	authv1.RegisterAuthServiceServer(server, authGRPCService)
	
	// Enable reflection for debugging with tools like grpcurl
	reflection.Register(server)
	
	return &GRPCServer{
		server: server,
		port:   port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	
	fmt.Printf("gRPC Server starting on port %d...\n", s.port)
	return s.server.Serve(lis)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}
