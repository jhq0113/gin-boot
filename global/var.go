package global

import (
	"gin-boot/conf"
	"gin-boot/mysql"
	"gin-boot/redis"
)

var (
	BootRedisGroup *redis.Group
	BootMysqlGroup *mysql.Group
)

func init() {
	BootRedisGroup = redis.NewGroup(conf.Conf.Redis["boot"])
	BootMysqlGroup = mysql.NewGroup(conf.Conf.Mysql["boot"])
}
