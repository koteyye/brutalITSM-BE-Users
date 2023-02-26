package handler

import (
	"brutalITSM-BE-Users/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		users := api.Group("/users")
		{
			users.GET("/", h.setRoleAdmin, h.checkRights, h.getUsers)
			users.GET("/:id", h.setRoleAdmin, h.checkRights, h.getUserById)
			users.POST("/create", h.setRoleAdmin, h.checkRights, h.createUser)
			users.DELETE("/delete/:id", h.setRoleAdmin, h.checkRights, h.deleteUser)
		}
	}

	return router
}
