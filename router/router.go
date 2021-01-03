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

func Load() *gin.Engine {
	boot := &controllers.Boot{}
	router.GET("/boot/client-ip", boot.ClientIp)
	redis := &controllers.Redis{}
	router.GET("/redis/group", redis.Group)
	router.GET("/redis/pipeline", redis.PipeLine)
	router.GET("/redis/multi", redis.Multi)

	mysql := &controllers.Mysql{}
	router.GET("/mysql/insert", mysql.Insert)
	router.GET("/mysql/batch-insert", mysql.BatchInsert)
	router.GET("/mysql/some", mysql.Some)
	router.GET("/mysql/one", mysql.One)
	router.GET("/mysql/delete", mysql.Delete)

	return router
}
