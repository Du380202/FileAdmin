package api

import (
	"backend/service"

	"github.com/gin-gonic/gin"
)

func UploadFileHandler(c *gin.Context) {
	service.UploadFile(c)
}

func DownloadFileHandler(c *gin.Context) {
	service.DownloadFile(c)
}

func GetFileHandler(c *gin.Context) {
	service.GetFile(c)
}

func SearchFileHandler(c *gin.Context) {
	service.SearchFile(c)
}
