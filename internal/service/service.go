package service

import (
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"github.com/minio/minio-go/v7"
)

type Authorization interface {
}

type User interface {
}

type Search interface {
}

type Service struct {
	Authorization
	User
	Search
}

func NewService(repos *postgres.Repository, s3 *minio.Client) *Service {
	return &Service{}
}
