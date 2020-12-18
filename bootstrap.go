package main

import (
	"gin-boot/conf"
	"gin-boot/redis"

	"github.com/gin-gonic/gin"
)

func Bootstrap() {
	gin.SetMode(conf.Conf.App.Mode)
	redis.GetPool(conf.Conf.Redis["boot"])
}
