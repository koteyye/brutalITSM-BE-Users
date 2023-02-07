package main

import (
	"brutalITSM-BE-Users/internal/user"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Println("Create router")
	router := httprouter.New()

	log.Println("register user handler")
	handler := user.NewHandler()
	handler.Register(router)

	start(router)

}

func start(router *httprouter.Router) {
	log.Println("start application")

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("server is listening port")
	log.Fatalln(server.Serve(listener))
}
