package router

import (
	"backend/api"
	"backend/middleware"
	"backend/setup"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterRouter(setup.R)
}

func RegisterRouter(r *gin.Engine) {
	scpGroup := r.Group("/scp")
	scpGroup.Use(middleware.JWTAuthMiddleware())
	{
		scpGroup.POST("/transfer", api.TransferFileHandler)
	}
}


