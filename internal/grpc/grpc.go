package grpc

import (
	"context"
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
)

type userServer struct {
	pb.UserServiceServer
}

func (s *userServer) CheckToken(ctx context.Context, req *pb.RequestToken) (*pb.RequestToken, error) {

}
