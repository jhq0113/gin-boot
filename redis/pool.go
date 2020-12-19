package redis

import (
	"fmt"
	"time"

	redigo "github.com/garyburd/redigo/redis"

	"gin-boot/conf"
)

type Pool struct {
	pool *redigo.Pool
}

func NewPool(option *conf.RedisOption) *Pool {
	return &Pool{
		pool: &redigo.Pool{
			MaxConnLifetime: time.Second * time.Duration(option.MaxConnLifetime),
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
		},
	}
}

func (this *Pool) Get() *Redis {
	return &Redis{
		conn: this.pool.Get(),
	}
}
