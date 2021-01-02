package mysql

import (
	"bytes"
	"fmt"
	"strings"
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

func (p *Pool) Find() *Query {
	return AcquireQuery()
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
	return p.db.QueryRow(sqlStr, args...)
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

func (p *Pool) BatchInsert(table string, rows []map[string]interface{}) (sql.Result, error) {
	fields := make([]string, 0, len(rows[0]))
	args := make([]interface{}, 0, len(rows)*len(rows[0]))
	values := make([]string, 0, len(rows))

	value := make([]byte, 0, 1+2*len(rows[0]))
	value = append(value, '(')
	for field, arg := range rows[0] {
		fields = append(fields, field)
		value = append(value, '?', ',')
		args = append(args, arg)
	}
	value[len(value)-1] = ')'

	values = append(values, string(value))

	for start := 1; start < len(rows); start++ {
		for _, arg := range rows[start] {
			args = append(args, arg)
		}
		values = append(values, string(value))
	}

	return p.db.Exec(fmt.Sprintf("INSERT INTO %s(%s)VALUES%s", table, strings.Join(fields, ","), strings.Join(values, ",")), args...)
}

func (p *Pool) UpdateAll(table string, set map[string]interface{}, where map[string]interface{}) (sql.Result, error) {
	sqlBuffer := bytes.NewBufferString(fmt.Sprintf("UPDATE %s SET ", table))

	args := utils.AcquireArgs()
	defer utils.ReleaseArgs(args)

	var num = 0
	for field, arg := range set {
		if num > 0 {
			sqlBuffer.WriteByte(',')
		} else {
			num++
		}
		sqlBuffer.Write([]byte(field))
		sqlBuffer.Write([]byte("=?"))
		args = append(args, arg)
	}

	if len(where) > 0 {
		condition, params := buildWhere(where)
		defer utils.ReleaseArgs(params)

		sqlBuffer.Write(condition)
		args = append(args, params...)
	}

	return p.db.Exec(sqlBuffer.String(), args...)
}

func (p *Pool) DeleteAll(table string, where map[string]interface{}) (sql.Result, error) {
	sqlBuffer := bytes.NewBufferString("DELETE FROM ")
	sqlBuffer.Write([]byte(table))

	if len(where) > 0 {
		condition, args := buildWhere(where)
		defer utils.ReleaseArgs(args)

		sqlBuffer.Write(condition)
		return p.db.Exec(sqlBuffer.String(), args...)
	}
	return p.db.Exec(sqlBuffer.String())
}
