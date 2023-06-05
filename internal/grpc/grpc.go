package grpc2

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
)

type GRPC struct {
	services *service.Service
	pb.UserServiceServer
}

type Error struct {
	Message string     `json:"message"`
	Status  codes.Code `json:"-"`
}

func NewGRPC(services *service.Service) *GRPC {
	return &GRPC{services: services}
}

func (e Error) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

func (e Error) GRPCStatus() *status.Status {
	return status.New(e.Status, e.Error())
}

func newError(status codes.Code, msg string) error {
	return &Error{
		Status:  status,
		Message: msg,
	}
}

func (s *GRPC) GetByToken(ctx context.Context, req *pb.RequestToken) (*pb.ResponseUser, error) {
	userId, err := s.services.ParseToken(req.Token)
	if err != nil {
		var errToken *jwt.ValidationError
		if errors.As(err, &errToken) {
			if errToken.Errors == 20 {
				return nil, newError(codes.Unauthenticated, err.Error())
			} else if errToken.Errors == 16 {
				return nil, newError(codes.PermissionDenied, err.Error())

			} else {
				return nil, newError(codes.Internal, err.Error())
			}
		} else {
			return nil, newError(codes.Internal, err.Error())
		}

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

func (s *GRPC) GetByUserList(ctx context.Context, req *pb.RequestUsers) (*pb.ResponseUsers, error) {
	userList, err := s.services.GetUserList(req.Id)
	if err != nil {
		return nil, err
	}
	var usersProto []*pb.User
	for _, users := range userList {
		usersProto = append(usersProto, &pb.User{
			Id:         users.Id,
			LastName:   users.Lastname,
			FirstName:  users.Firstname,
			SurName:    users.Surname,
			MimeType:   users.MimeType,
			BucketName: users.BucketName,
			FileName:   users.FileName,
		})
	}
	return &pb.ResponseUsers{UserList: usersProto}, nil
}
