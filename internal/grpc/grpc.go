package grpc

import (
	"context"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
)

type GRPC struct {
	service *service.Service
	pb.UserServiceServer
}

func (s *GRPC) CheckToken(ctx context.Context, req *pb.RequestToken) (*pb.ResponseUser, error) {
	user, err := s.service.Me(req)
	if err != nil {
		return nil, err
	}
	return &pb.ResponseUser{
		Id:               user.Id,
		Login:            user.Login,
		LastName:         user.Lastname,
		FirstName:        user.Firstname,
		SurName:          user.Surname,
		Job:              user.Job,
		Org:              user.Org,
		Roles:            user.RolesList,
		AvatarMimeType:   user.Avatar.MimeType,
		AvatarBucketName: user.Avatar.BucketName,
		AvatarFileName:   user.Avatar.FileName,
	}, nil
}
