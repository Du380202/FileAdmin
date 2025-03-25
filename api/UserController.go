package api

import (
	"backend/service"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	service.Register(c)
}

func LoginHandler(c *gin.Context) {
	service.Login(c)
}
