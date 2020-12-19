package redis

import (
	"hash/crc32"

	"gin-boot/conf"
)

type Group struct {
	pool []*Pool
	size uint32
}

func NewGroup(options []conf.RedisOption) *Group {
	poolSize := len(options)
	if poolSize < 1 {
		panic("redis options长度小于1")
	}

	group := &Group{
		size: uint32(poolSize),
		pool: make([]*Pool, 0, poolSize),
	}

	for index := 0; index < poolSize; index++ {
		group.pool = append(group.pool, NewPool(&options[index]))
	}

	return group
}

func (this *Group) GetPool(key []byte) *Pool {
	if this.size < 2 {
		return this.pool[0]
	}
	index := crc32.ChecksumIEEE(key) % this.size
	return this.pool[index]
}
