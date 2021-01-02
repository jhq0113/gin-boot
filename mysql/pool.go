package mysql

import (
	"bytes"
	"fmt"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"gin-boot/conf"
	"gin-boot/utils"
)

type Pool struct {
	db *sql.DB
}

func NewPool(option *conf.MysqlOption) (*Pool, error) {
	db, err := sql.Open("mysql", option.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(option.MaxConnLifetime) * time.Second)
	db.SetMaxIdleConns(option.MaxIdleConns)
	db.SetMaxOpenConns(option.MaxOpenConns)

	return &Pool{
		db: db,
	}, nil
}

func (p *Pool) Db() *sql.DB {
	return p.db
}

func (p *Pool) Query(sqlStr string, args ...interface{}) (*sql.Rows, error) {
	return p.db.Query(sqlStr, args...)
}

func (p *Pool) Execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(sqlStr, args...)
}

func (p *Pool) All(query *Query) (*sql.Rows, error) {
	sqlStr, args := BuildQuery(query)
	defer func() {
		ReleaseQuery(query)
		utils.ReleaseArgs(args)
	}()
	return p.db.Query(sqlStr, args...)
}

func (p *Pool) One(query *Query) *sql.Row {
	query.limit = 1
	sqlStr, args := BuildQuery(query)
	defer func() {
		ReleaseQuery(query)
		utils.ReleaseArgs(args)
	}()
	return p.db.QueryRow(sqlStr, args)
}

func (p *Pool) Insert(table string, columns map[string]interface{}) (sql.Result, error) {
	sqlBuffer := bytes.NewBufferString(fmt.Sprintf("INSERT INTO %s(", table))

	args := utils.AcquireArgs()
	defer utils.ReleaseArgs(args)

	values := make([]byte, 0, 7+2*len(columns))
	values = append(values, []byte("VALUES(")...)

	for field, arg := range columns {
		if len(values) > 7 {
			sqlBuffer.WriteByte(',')
		}
		sqlBuffer.Write([]byte(field))

		values = append(values, '?', ',')
		args = append(args, arg)
	}
	sqlBuffer.WriteByte(')')
	values[len(values)-1] = ')'
	sqlBuffer.Write(values)

	return p.db.Exec(sqlBuffer.String(), args...)
}
