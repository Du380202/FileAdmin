package router

import (
	"backend/api"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (user *UserRouter) RegisterRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", api.RegisterHandler)
		userGroup.POST("/login", api.LoginHandler)

	}
}
