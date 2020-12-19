package global

import (
	"gin-boot/conf"
	"gin-boot/redis"
)

var (
	BootRedisGroup *redis.Group
)

func init() {
	BootRedisGroup = redis.NewGroup(conf.Conf.Redis["boot"])
}
