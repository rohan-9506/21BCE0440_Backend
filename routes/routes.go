package routes

import (
	"file-sharing-system/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Add rate limit middleware to the router
	r.Use(api.RateLimitMiddleware())

	r.POST("/api/register", api.RegisterHandler)
	r.POST("/api/login", api.LoginHandler)
	r.POST("/api/upload", api.UploadHandler)

	// WebSocket route
	r.GET("/ws", api.WebSocketHandler)

	r.GET("/api/files", api.GetFilesHandler)
	r.GET("/api/share/:file_id", api.ShareFileHandler)

	authRoutes := r.Group("/api")
	authRoutes.Use(api.AuthMiddleware())
	{
		//pvt routes here which require authentication like below
		//authRoutes.GET("/api/files", api.GetFilesHandler)
		//authRoutes.GET("/api/share/:file_id", api.ShareFileHandler)
		//authRoutes.POST("/api/upload", api.UploadHandler)
		//authRoutes.GET("/api/ws", api.WebSocketHandler)
		//authRoutes.POST("/api/logout", api.LogoutHandler)
	}

	return r
}
