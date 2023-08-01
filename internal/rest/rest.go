package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

type Rest struct {
	services *service.Service
}

func NewRest(services *service.Service) *Rest {
	return &Rest{services: services}
}

func (h *Rest) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"x-requested-with, Content-Type, origin, authorization, accept, x-access-token"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.GET("/me", h.userIdentity, h.me)
	}

	api := router.Group("/api", h.userIdentity)
	{
		users := api.Group("/users")
		{
			users.GET("/", h.setRoleAdmin, h.checkRights, h.getUsers)
			users.GET("/:id", h.setRoleAdmin, h.checkRights, h.getUserById)
			users.POST("/create", h.setRoleAdmin, h.checkRights, h.createUser)
			users.DELETE("/delete/:id", h.setRoleAdmin, h.checkRights, h.deleteUser)
			users.POST("/avatar/upload/:id", h.setRoleAdmin, h.checkRights, h.uploadFile)
			users.GET("/roles", h.setRoleAdmin, h.checkRights, h.getRoles)
		}
		search := api.Group("/search")
		{
			search.GET("/job/:jobName", h.setRoleAdmin, h.checkRights, h.searchJob)
			search.GET("/org/:orgName", h.setRoleAdmin, h.checkRights, h.searchOrg)
		}
		settings := api.Group("/settings")
		{
			settings.POST("/add", h.setRoleAdmin, h.checkRights)
			settings.DELETE("/delete", h.setRoleAdmin, h.checkRights)
			settings.PUT("/edit", h.setRoleAdmin, h.checkRights)
		}
	}

	return router
}
