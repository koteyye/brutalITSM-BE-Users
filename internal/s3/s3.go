package s3

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func InitMinioClient() (*minio.Client, error) {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initialization s3 configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env s3 variables %s", err.Error())
	}

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

	return minioClient, err
}

func initConfig() error {
	viper.AddConfigPath("server")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
