package redis

import (
	"fmt"
	"gin-boot/global"

	redigo "github.com/garyburd/redigo/redis"
)

func init() {

}

type Pool struct {
	pool     []*redigo.Pool
	poolSize int
}

func (p *Pool) add(option *global.RedisOption) {
	p.pool = append(p.pool, &redigo.Pool{
		Dial: func () (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", fmt.Sprintf("%s:%s",option.Host, option.Port))
			 if err != nil {
				return nil, err
			 }

			 if len(option.Auth) > 0 {
				 if _, err := c.Do("AUTH", option.Auth); err != nil {
					 _ = c.Close()
					 return nil, err
				 }
			 }

			 if _, err := c.Do("SELECT", option.Db); err != nil {
				 _ = c.Close()
				return nil, err
			 }
			 return c, nil
		},
	})
}

func GetPool(options []global.RedisOption) *Pool {
	poolSize := len(options)
	if poolSize < 1 {
		panic("redis options长度小于1")
	}

	pool := &Pool{
		poolSize: poolSize,
		pool:     make([]*redigo.Pool, 0, poolSize),
	}

	for index := 0; index < poolSize; index++ {
		pool.add(&options[index])
	}

	return pool
}
