package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Controller struct{}

func (this *Controller) clientIp(ctx *gin.Context) string {
	return ctx.ClientIP()
}

func QueryString(ctx *gin.Context, key string, defaultValue string) string {
	argValue, exists := ctx.GetQuery(key)
	if !exists {
		return defaultValue
	}
	return argValue
}

func QueryInt(ctx *gin.Context, key string, defaultValue int) int {
	argValue, exists := ctx.GetQuery(key)
	if !exists {
		return defaultValue
	}

	convertValue, err := strconv.Atoi(argValue)
	if err != nil {
		return defaultValue
	}

	return convertValue
}

func QueryInt64(ctx *gin.Context, key string, defaultValue int64) int64 {
	argValue, exists := ctx.GetQuery(key)
	if !exists {
		return defaultValue
	}

	convertValue, err := strconv.ParseInt(argValue, 10, 64)
	if err != nil {
		return defaultValue
	}

	return convertValue
}
