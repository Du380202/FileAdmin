package api

import (
	"backend/service"

	"github.com/gin-gonic/gin"
)

func TransferFileHandler(c *gin.Context) {
	service.TransferFile(c)
}
