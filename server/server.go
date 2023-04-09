package server

import (
	pb "github.com/koteyye/brutalITSM-BE-Users/proto"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

func (s *Server) Run(port string, handler http.Handler) error {

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	s.grpcServer = &grpc.Server{
		pb
	}

	return s.httpServer.ListenAndServe()
}

