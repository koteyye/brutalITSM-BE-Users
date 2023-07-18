package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/koteyye/brutalITSM-BE-Users/config"
	grpc2 "github.com/koteyye/brutalITSM-BE-Users/internal/grpc"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"github.com/koteyye/brutalITSM-BE-Users/internal/rest"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
	brutalitsm "github.com/koteyye/brutalITSM-BE-Users/server"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title User Service API
// @version 1.0
// @description API Server for User Service BrutalITSM

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables %s", err.Error())
	}

	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("Get config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	// Postgres Client
	db, err := newPostgres(ctx, cfg.Storages.Postgres)
	if err != nil {
		logrus.Fatalf("Can't get postgres client pool: %v", err)
	}

	// Minio Client
	endpoint := cfg.Storages.Minio.URL
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

	go runGrpcServer(grpcHandler, cfg.Server.GRPC)

	restServer := new(brutalitsm.Server)
	if err := restServer.Run(cfg.Server.HTTP.Listen, handler.InitRoutes()); err != nil {
		logrus.Fatalf("Error occuped while runing Rest server :%s", err.Error())
	}

}

func runGrpcServer(grpcHandler *grpc2.GRPC, cfg *config.GrpcServerConfig) {
	lis, err := net.Listen("tcp", cfg.Listen)
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

func newPostgres(ctx context.Context, cfg *config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("can't create db: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = db.PingContext(dbCtx)
	if err != nil {
		return nil, fmt.Errorf("can't ping db: %w", err)
	}

	logrus.Infoln("Migration started")
	m, err := migrate.New("file://db/migrations",
		cfg.DSN)
	if err != nil {
		logrus.Fatalf("Migrate error: %v", err)
	}
	if err := m.Up(); err != nil {
		switch err {
		case nil:
			break
		case migrate.ErrNoChange:
			logrus.Info("Migrate no change")
			return db, nil
		default:
			logrus.Fatalf("Migrate error: %v", err)
			return db, err
		}
	}
	return db, nil
}
