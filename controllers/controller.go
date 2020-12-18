package controllers

import "github.com/gin-gonic/gin"

type Controller struct {}

func (this *Controller) clientIp(ctx *gin.Context) string {
	return ctx.ClientIP()
}
