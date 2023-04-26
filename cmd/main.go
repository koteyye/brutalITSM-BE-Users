package main

import (
	"github.com/joho/godotenv"
	grpc2 "github.com/koteyye/brutalITSM-BE-Users/internal/grpc"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"github.com/koteyye/brutalITSM-BE-Users/internal/rest"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
	brutalitsm "github.com/koteyye/brutalITSM-BE-Users/server"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables %s", err.Error())
	}
	// Postgres Client
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("faild to initialize db: %s", err.Error())
	}

	// Minio Client
	endpoint := viper.GetString("minio.url")
	accessKeyId := os.Getenv("KEY_ID")
	secretAccessKey := os.Getenv("SECRET_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Printf("%#v\n", minioClient) // minio is now set up

	//init internal
	repos := postgres.NewRepository(db)
	services := service.NewService(repos, minioClient)
	handler := rest.NewRest(services)
	grpcHandler := grpc2.NewGRPC(services)

	go runGrpcServer(grpcHandler)

	restServer := new(brutalitsm.Server)
	if err := restServer.Run(viper.GetString("restPort"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("Error occuped while runing Rest server :%s", err.Error())
	}

}

func runGrpcServer(grpcHandler *grpc2.GRPC) {
	lis, err := net.Listen("tcp", viper.GetString("grpcPort"))
	if err != nil {
		logrus.Fatalf("Failed to start GRPC server \n %v \n", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcHandler)
	logrus.Info("GRPC Server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("Failed to start GRPC: \n %v \n", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("server")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
