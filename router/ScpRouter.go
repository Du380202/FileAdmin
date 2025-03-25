package router

import (
	"backend/api"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

type SCPRouter struct{}

func (scp *SCPRouter) RegisterRouter(r *gin.Engine) {
	scpGroup := r.Group("/scp")
	scpGroup.Use(middleware.JWTAuthMiddleware())
	{
		scpGroup.POST("/transfer", api.TransferFileHandler)
	}
}
