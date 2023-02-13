package main

import (
	"brutalITSM-BE-Users/internal/config"
	"brutalITSM-BE-Users/internal/user"
	user2 "brutalITSM-BE-Users/internal/user/db"
	"brutalITSM-BE-Users/pkg/client/postgresql"
	"brutalITSM-BE-Users/pkg/logging"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	repository := user2.NewRepository(postgreSQLClient, logger)

	newUser := user.User{
		Login:    "Egorka",
		Password: "12345",
	}
	err = repository.Create(context.TODO(), &newUser)
	if err != nil {
		logger.Fatalf("%v", err)
	} else {
		newPerson := user.Person{
			FirstName:  "Егорка",
			LastName:   "Аналитиков",
			MiddleName: "",
			JobName:    "Аналитик",
			OrgName:    "Анал-литическая",
			UserId:     newUser.ID,
		}
		err = repository.CreatePerson(context.TODO(), &newPerson)
		if err != nil {
			logger.Fatalf("%v", err)
		}
		logger.Infof("NewUser: %v NewPerson: %v", newUser, newPerson)
	}

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
