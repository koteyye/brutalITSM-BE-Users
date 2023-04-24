package server

import (
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":8081"
)

type userServer struct {
	pb.UserServiceServer
}

func newGrpcServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	return grpcServer
}

func RunGrpcSrv() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("Failed to start the server %v\n", err)
	}
	grpcServer := newGrpcServer()
	pb.RegisterUserServiceServer(grpcServer, &userServer{})
	log.Printf("Server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("Failed to start: %v\n", err)
	}

	return grpcServer.Serve(lis)
}

func Stop(s *grpc.Server) {
	s.Stop()
}
