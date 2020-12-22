package mysql

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

const (
	formatIn = " (`%s` IN (%s)) AND"
	formatBetween = " (`%s` BETWEEN ? AND ?) AND"
	formatLike = " (`%s` LIKE ?) AND"
	formatCom = " (`%s` %s ?) AND"
	inHolder =  "?,"
	defaultLimit  = 1000
)

var (
	defaultColumns       = []string{"*"}
	wherePrefix          = []byte(" WHERE ")
	queryPool            = &sync.Pool{
		New: func() interface{} {
			query := &Query{
				limit:   defaultLimit,
				columns: defaultColumns,
			}

			return query
		},
	}

	argsPool = &sync.Pool{
		New: func() interface{} {
			return make([]interface{}, 0, 8)
		},
	}
)

type Query struct {
	table   string
	columns []string
	where   map[string]interface{}
	group   []string
	order   []string
	offset  int64
	limit   int64
}

//---------------------查询对象池--------------------------

func AcquireQuery() *Query {
	return queryPool.Get().(*Query)
}

func ReleaseQuery(query *Query) {
	query = query.reset()
	queryPool.Put(query)
}

//---------------------查询对象池--------------------------

//---------------------参数池--------------------------

func acquireArgs() []interface{} {
	return argsPool.Get().([]interface{})
}

func releaseArgs(args []interface{}) {
	args = args[:0]
	argsPool.Put(args)
}

//---------------------参数池--------------------------

func (q *Query) reset() *Query {
	q.table = ""
	q.columns = defaultColumns
	q.offset = 0
	q.limit = defaultLimit
	q.where = nil
	q.group = nil
	q.order = nil

	return q
}

func (q *Query) From(table string) *Query {
	q.table = table
	return q
}

func (q *Query) Select(columns ...string) *Query {
	q.columns = columns
	return q
}

func (q *Query) Where(where map[string]interface{}) *Query {
	q.where = where
	return q
}

func (q *Query) Group(group ...string) *Query {
	q.group = group
	return q
}

func (q *Query) Order(order ...string) *Query {
	q.order = order
	return q
}

func (q *Query) Offset(offset int64) *Query {
	q.offset = offset
	return q
}

func (q *Query) Limit(offset int64, limit int64) *Query {
	if limit == 0 {
		limit = defaultLimit
	}
	q.limit = limit
	q.offset = offset
	return q
}

func buildWhere(where map[string]interface{}, args []interface{}) (condition []byte) {
	if len(where) < 1 {
		return
	}

	buf := bytes.NewBuffer(wherePrefix)

	var (
		operator string
		position int
		inLength int
	)

	for field, value := range where {
		operator = "="
		position = strings.Index(field, " ")
		if position > 0 {
			operator = field[position+1:]
			field = field[:position]
		}

		if val,ok := value.([]interface{}); ok {
			operator = "IN"
			args = append(args, val...)
			inLength = len(val)
		} else {
			args = append(args, value)
		}

		switch operator {
		case "IN":
			buf.WriteString(fmt.Sprintf(formatIn, field,  strings.Repeat(inHolder, inLength)[:2*inLength-1]))
		case "BETWEEN":
			buf.WriteString(fmt.Sprintf(formatBetween, field))
		case "LIKE":
			buf.WriteString(fmt.Sprintf(formatLike, field))
		default:
			buf.WriteString(fmt.Sprintf(formatCom, field, operator))
		}
	}

	whereCon := buf.Bytes()
	return whereCon[:len(whereCon) - 3]
}

func BuildQuery(query *Query) (sql string, args []interface{}) {
	buf := bytes.NewBufferString("SELECT")
}
