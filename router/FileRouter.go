package router

import (
	"backend/api"

	"github.com/gin-gonic/gin"
)

type FileRouter struct{}

func (s *FileRouter) RegisterRouter(r *gin.Engine) {
	fileGroup := r.Group("/file")
	{
		fileGroup.GET("/search", api.SearchFileHandler)
		fileGroup.GET("/download/:path/:filename", api.DownloadFileHandler)
		fileGroup.POST("/upload", api.UploadFileHandler)
	}
}
