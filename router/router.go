package router

import (
	"backend/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"}, // Cho phép frontend truy cập API
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/file/:path", api.UploadFileHandler)
	r.GET("/download/:filename", api.DownloadFileHandler)
	r.GET("/file", api.GetFileHandler)
	r.GET("/search", api.SearchFileHandler)
	r.POST("/transfer", api.TransferFileHandler)
	return r
}
