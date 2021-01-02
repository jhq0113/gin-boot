package mysql

import (
	"database/sql"
	"math/rand"

	"gin-boot/conf"
)

type Group struct {
	Masters []*Pool
	Slaves  []*Pool

	masterLen int
	slaveLen  int
}

func NewGroup(groupOption conf.MysqlGroupOption) *Group {
	if len(groupOption.Slaves) == 0 {
		groupOption.Slaves = groupOption.Masters
	}

	group := &Group{
		masterLen: len(groupOption.Masters),
		slaveLen:  len(groupOption.Slaves),
	}

	group.Masters = make([]*Pool, 0, group.masterLen)
	group.Slaves = make([]*Pool, 0, group.slaveLen)

	for index, _ := range groupOption.Masters {
		pool, err := NewPool(&groupOption.Masters[index])
		if err != nil {
			panic(err.Error())
		}

		group.Masters = append(group.Masters, pool)
	}

	for index, _ := range groupOption.Slaves {
		pool, err := NewPool(&groupOption.Slaves[index])
		if err != nil {
			panic(err.Error())
		}

		group.Slaves = append(group.Slaves, pool)
	}

	return group
}

func (g *Group) SelectPool(useMaster bool) *Pool {
	if useMaster {
		if g.masterLen == 1 {
			return g.Masters[0]
		}

		return g.Masters[rand.Intn(g.masterLen)]
	}

	if g.slaveLen == 1 {
		return g.Slaves[0]
	}

	return g.Slaves[rand.Intn(g.slaveLen)]
}

func (g *Group) Query(sqlStr string, args []interface{}, useMaster bool) (*sql.Rows, error) {
	return g.SelectPool(useMaster).Query(sqlStr, args)
}

func (g *Group) Execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return g.SelectPool(true).Execute(sqlStr, args...)
}
