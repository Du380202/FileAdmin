package router

import (
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(routers ...Router) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	for _, router := range routers {
		router.RegisterRouter(r)
	}

	return r
}
