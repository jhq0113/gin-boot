package redis

import (
	"fmt"
	"gin-boot/conf"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

type Pool struct {
	pool     []*redigo.Pool
	poolSize int
}

func (p *Pool) add(option *conf.RedisOption) {
	pool := &redigo.Pool{
		MaxConnLifetime: option.MaxConnLifetime,
		MaxIdle:         option.MaxIdle,
		MaxActive:       option.MaxActive,
		Wait:            option.Wait,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp",
				fmt.Sprintf("%s:%s", option.Host, option.Port),
				redigo.DialConnectTimeout(time.Millisecond*time.Duration(option.ConnectTimeout)),
				redigo.DialReadTimeout(time.Millisecond*time.Duration(option.ReadTimeout)),
				redigo.DialReadTimeout(time.Millisecond*time.Duration(option.ReadTimeout)),
			)
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
	}

	p.pool = append(p.pool, pool)
}

func (this *Pool) Get(key []byte) []byte {

}

func GetPool(options []conf.RedisOption) *Pool {
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
