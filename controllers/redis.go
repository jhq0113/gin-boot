package controllers

import (
	"bytes"
	"fmt"
	"gin-boot/conf"
	"gin-boot/global"
	"gin-boot/redis"

	"github.com/gin-gonic/gin"
)

type Redis struct {
	Controller
}

func (this *Redis) PipeLine(ctx *gin.Context) {
	pool := redis.NewPool(&conf.Conf.Redis["boot"][0])
	redis := pool.Get()
	defer redis.Close()

	multi := redis.Multi()
	result, err := multi.Send("SET", "boot:pipeline", "pipeline:boot").
		Send("GET", "boot:pipeline").
		ExecPipeline(redis)
	fmt.Println(result, err)
	ctx.Writer.WriteString("ok")
}

func (this *Redis) Multi(ctx *gin.Context) {
	pool := redis.NewPool(&conf.Conf.Redis["boot"][0])
	redis := pool.Get()
	defer redis.Close()

	multi := redis.Multi()
	result, _ := multi.Send("SET", "boot:multi", "multi:boot").
		Send("GET", "boot:multi").
		ExecMulti(redis)

	buf := bytes.NewBuffer(nil)
	for _, reply := range result {

	}
	ctx.Writer.WriteString("ok")
}

func (this *Redis) Group(ctx *gin.Context) {
	key := []byte("boot:redis")
	redis := global.BootRedisGroup.GetPool(key).Get()
	//释放连接到连接池
	defer redis.Close()

	redis.SetTimeout(key, []byte("redis:group:test"), 60)
	val, _ := redis.Get(key)
	ctx.Writer.Write(val)
}
