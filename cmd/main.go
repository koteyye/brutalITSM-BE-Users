package main

import (
	brutalitsm "brutalITSM-BE-Users"
	"brutalITSM-BE-Users/pkg/handler"
	"brutalITSM-BE-Users/pkg/repository"
	"brutalITSM-BE-Users/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
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

	//minio
	endpoint := viper.GetString("minio.url")
	accessKeyId := os.Getenv("KEY_ID")
	secretAccessKey := os.Getenv("SECRET_KEY")
	useSSL := false

	// Init minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Printf("%#v\n", minioClient) // minio is now set up

	repos := repository.NewRepository(db)
	services := service.NewService(repos, minioClient)
	handlers := handler.NewHandler(services)

	srv := new(brutalitsm.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occuped while runing http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
