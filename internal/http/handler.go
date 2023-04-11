package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	"time"
)

type Http struct {
	services *service.Service
}

func NewHttp(services *service.Service) *Http {
	return &Http{services: services}
}

func (h *Http) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"x-requested-with, Content-Type, origin, authorization, accept, x-access-token"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return router
}
