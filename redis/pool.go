package redis

import (
	"fmt"
	"hash/crc32"
	"strings"
	"sync"
	"time"

	redigo "github.com/garyburd/redigo/redis"

	"gin-boot/conf"
)

var (
	argsPool = &sync.Pool{
		New: func() interface{} {
			return make([]interface{}, 0, 8)
		},
	}
)

func acquireArgs() []interface{} {
	return argsPool.Get().([]interface{})
}

func releaseArgs(args []interface{}) {
	args = args[:0]
	argsPool.Put(args)
}

type Pool struct {
	pool     []*redigo.Pool
	poolSize uint32
}

func (p *Pool) add(option *conf.RedisOption) {
	pool := &redigo.Pool{
		MaxConnLifetime: time.Duration(option.MaxConnLifetime),
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

func (this *Pool) getConnection(key []byte) redigo.Conn {
	if this.poolSize < 2 {
		return this.pool[0].Get()
	}

	index := crc32.ChecksumIEEE(key) % this.poolSize
	return this.pool[index].Get()
}

func (this *Pool) Do(cmd string, key []byte, params ...interface{}) ([]byte, error) {
	args := acquireArgs()
	args = append(args, params)
	return redigo.Bytes(this.getConnection(key).Do(cmd, args...))
}

//----------------------------------String----------------------------------------------
func (this *Pool) Get(key []byte) ([]byte, error) {
	return redigo.Bytes(this.getConnection(key).Do("GET", key))
}

func (this *Pool) Set(key []byte, value []byte, params ...interface{}) bool {
	args := acquireArgs()
	args = append(args, key, value)
	args = append(args, params...)
	receive, _ := redigo.String(this.getConnection(key).Do("SET", args...))
	releaseArgs(args)
	return strings.ToUpper(receive) == "OK"
}

//----------------------------------String----------------------------------------------

func GetPool(options []conf.RedisOption) *Pool {
	poolSize := len(options)
	if poolSize < 1 {
		panic("redis options长度小于1")
	}

	pool := &Pool{
		poolSize: uint32(poolSize),
		pool:     make([]*redigo.Pool, 0, poolSize),
	}

	for index := 0; index < poolSize; index++ {
		pool.add(&options[index])
	}

	return pool
}
