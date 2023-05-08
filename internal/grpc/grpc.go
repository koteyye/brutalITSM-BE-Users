package grpc2

import (
	"context"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
)

type GRPC struct {
	services *service.Service
	pb.UserServiceServer
}

func NewGRPC(services *service.Service) *GRPC {
	return &GRPC{services: services}
}

func (s *GRPC) GetByToken(ctx context.Context, req *pb.RequestToken) (*pb.ResponseUser, error) {
	userId, err := s.services.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	user, err := s.services.Me(userId)
	if err != nil {
		return nil, err
	}
	return &pb.ResponseUser{
		Id:          user.Id,
		Login:       user.Login,
		LastName:    user.Lastname,
		FirstName:   user.Firstname,
		SurName:     user.Surname,
		Job:         user.Job,
		Org:         user.Org,
		Roles:       user.RolesList,
		Permissions: user.Permissions,
	}, nil
}

func (s *GRPC) GetByUserId(ctx context.Context, req *pb.RequestUser) (*pb.ResponseUser, error) {
	user, err := s.services.GetUserById(req.UserId)
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
		Permissions:      user.Permissions,
		AvatarMimeType:   user.Avatar.MimeType,
		AvatarBucketName: user.Avatar.BucketName,
		AvatarFileName:   user.Avatar.FileName,
	}, nil
}

//func (s *GRPC) GetByUsersId(ctx context.Context, req *pb.RequestUsers) (*pb.ResponseShortUsers, error) {
//	userList, err := s.services.GetUserList(req.Id)
//	if err != nil {
//		return nil, err
//	}
//	return &pb.ResponseShortUsers{
//		Id: userList.
//	}, nil
//}
