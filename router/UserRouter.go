package router

import (
	"backend/api"
	"backend/setup"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterUserRouter(setup.R)
}

func RegisterUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", api.RegisterHandler)
		userGroup.POST("/login", api.LoginHandler)

	}
}
