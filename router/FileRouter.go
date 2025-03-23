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
		fileGroup.POST("/download", api.DownloadFileHandler)
		fileGroup.POST("/upload", api.UploadFileHandler)
	}
}
