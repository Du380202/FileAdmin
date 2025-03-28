package router

import (
	"backend/api"
	"backend/setup"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterFileRouter(setup.R)
}

func RegisterFileRouter(r *gin.Engine) {
	fileGroup := r.Group("/file")
	{
		fileGroup.GET("/", api.GetFileHandler)
		fileGroup.GET("/search", api.SearchFileHandler)
		fileGroup.GET("/fulltext", api.FullTextSearchHandler)
		fileGroup.GET("/download/:path/:filename", api.DownloadFileHandler)
		fileGroup.POST("/upload", api.UploadFileHandler)
	}
}
