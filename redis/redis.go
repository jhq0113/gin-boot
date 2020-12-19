package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"strings"
	"sync"
)

var (
	//参数池
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

type Redis struct {
	conn redigo.Conn
}

//释放连接到连接池，并非真正的close
func (this *Redis) Close() error {
	return this.conn.Close()
}

func (this *Redis) Do(cmd string, args ...interface{}) (interface{}, error) {
	return this.conn.Do(cmd, args...)
}

func (this *Redis) Send(cmd string, args ...interface{}) error {
	return this.conn.Send(cmd, args...)
}

func (this *Redis) Multi() Multi {
	return acquireMulti()
}

//----------------------------------String----------------------------------------------
func (this *Redis) Get(key []byte) ([]byte, error) {
	return redigo.Bytes(this.conn.Do("GET", key))
}

func (this *Redis) Set(key []byte, params ...interface{}) bool {
	args := acquireArgs()
	args = append(args, key)
	args = append(args, params...)
	receive, _ := redigo.String(this.conn.Do("SET", args...))
	releaseArgs(args)
	return strings.ToUpper(receive) == "OK"
}

func (this *Redis) SetTimeout(key []byte, value interface{}, timeoutSecond int64) bool {
	receive, _ := redigo.String(this.conn.Do("SETEX", key, timeoutSecond, value))
	return strings.ToUpper(receive) == "OK"
}

//----------------------------------String----------------------------------------------
