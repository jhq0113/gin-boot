package mysql

import (
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"gin-boot/conf"
)

type Pool struct {
	db *sql.DB
}

func NewPool(option *conf.MysqlOption) (*Pool, error){
	db, err := sql.Open("mysql", option.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(option.MaxConnLifetime) * time.Second)
	db.SetMaxIdleConns(option.MaxIdleConns)
	db.SetMaxOpenConns(option.MaxOpenConns)

	return &Pool{
		db:db,
	}, nil
}

func (p *Pool) Db() *sql.DB {
	return p.db
}
