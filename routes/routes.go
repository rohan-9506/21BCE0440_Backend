package routes

import (
	"file-sharing-system/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/api/register", api.RegisterHandler)
	r.POST("/api/login", api.LoginHandler)

	authRoutes := r.Group("/api")
	authRoutes.Use(api.AuthMiddleware())
	{
		authRoutes.POST("/upload", api.UploadHandler)
	}

	return r
}
