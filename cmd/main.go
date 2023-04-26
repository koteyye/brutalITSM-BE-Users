package main

import (
	"github.com/joho/godotenv"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"github.com/koteyye/brutalITSM-BE-Users/internal/rest"
	"github.com/koteyye/brutalITSM-BE-Users/internal/s3"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
	brutalitsm "github.com/koteyye/brutalITSM-BE-Users/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

type GrpcUserServer struct {
	pb.UserServiceServer
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables %s", err.Error())
	}
	// Postgres Client
	db, _ := postgres.InitPostgres()
	// Minio Client
	minio, _ := s3.InitMinioClient()

	//init internal
	repos := postgres.NewRepository(db)
	services := service.NewService(repos, minio)
	handler := rest.NewRest(services)

	go runGrpcServer()

	go runRestServer(handler)

}

func runGrpcServer() {
	lis, err := net.Listen("tcp", viper.GetString("grpcPort"))
	if err != nil {
		logrus.Fatalf("Failed to start GRPC server \n %v \n", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &GrpcUserServer{})
	logrus.Info("GRPC Server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("Failed to start GRPC: \n %v \n", err)
	}
}

func runRestServer(handler *rest.Rest) {
	restServer := new(brutalitsm.Server)
	if err := restServer.Run(viper.GetString("restPort"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("Error occuped while runing Rest server :%s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("server")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
