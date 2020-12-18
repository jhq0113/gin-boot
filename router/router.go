package router

import (
	"gin-boot/controllers"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func Load() *gin.Engine{
	boot := &controllers.Boot{}
	router.GET("/boot/client-ip", boot.ClientIp)

	return router
}