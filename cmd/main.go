package main

import (
	"github.com/joho/godotenv"
	http "github.com/koteyye/brutalITSM-BE-Users/internal/http"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"github.com/koteyye/brutalITSM-BE-Users/internal/s3"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	brutalitsm "github.com/koteyye/brutalITSM-BE-Users/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig; err != nil {
		logrus.Fatalf("Error initialization postgres configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env postgres variables %s", err.Error())
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
	minio, _ := s3.InitMinioClient()

	//init internal
	repos := postgres.NewRepositry(db)
	services := service.NewService(repos, minio)
	handler := http.NewHttp(services)

	//run http server
	httpSrv := new(brutalitsm.Server)
	if err := httpSrv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("error occuped while runing http server: %s", err.Error())
	}

	//run gRPC server
	_ = brutalitsm.RunGrpcSrv()

}

func initConfig() error {
	viper.AddConfigPath("server")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
