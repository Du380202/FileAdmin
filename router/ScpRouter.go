package router

import (
	"backend/api"

	"github.com/gin-gonic/gin"
)

type SCPRouter struct{}

func (scp *SCPRouter) RegisterRouter(r *gin.Engine) {
	scpGroup := r.Group("/scp")
	{
		scpGroup.POST("/transfer", api.TransferFileHandler)
	}
}
