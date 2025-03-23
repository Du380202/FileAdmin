package router

import "github.com/gin-gonic/gin"

type Router interface {
	RegisterRouter(r *gin.Engine)
}
