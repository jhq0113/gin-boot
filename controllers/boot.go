package controllers

import (
	"github.com/gin-gonic/gin"
)

type Boot struct {
	Controller
}

func (this *Boot) ClientIp(ctx *gin.Context) {
	ctx.Writer.Write([]byte(this.clientIp(ctx)))
}
