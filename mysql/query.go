package mysql

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"gin-boot/utils"
)

const (
	formatIn = " (`%s` IN (%s)) AND"
	formatBetween = " (`%s` BETWEEN ? AND ?) AND"
	formatLike = " (`%s` LIKE ?) AND"
	formatCom = " (`%s` %s ?) AND"
	//SELECT fields FROM table WHERE condition GROUP BY fields HAVING aggregation ORDER BY fields LIMIT offset,limit
	selectSql = "SELECT %s FROM %s %s %s %s %s LIMIT %d,%d"
	inHolder =  "?,"
	defaultLimit  = 1000
)

var (
	defaultColumns       = "*"
	wherePrefix          = []byte("WHERE ")
	queryPool            = &sync.Pool{
		New: func() interface{} {
			query := &Query{
				limit:   defaultLimit,
				columns: defaultColumns,
			}

			return query
		},
	}
)

type Query struct {
	table   string
	columns string
	where   map[string]interface{}
	group   string
	having  string
	order   string
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

func (q *Query) reset() *Query {
	q.table = ""
	q.columns = defaultColumns
	q.offset = 0
	q.limit = defaultLimit
	q.where = nil
	q.group = ""
	q.having = ""
	q.order = ""

	return q
}

func (q *Query) From(table string) *Query {
	q.table = table
	return q
}

func (q *Query) Select(columns...string) *Query {
	q.columns = strings.Join(columns, ",")
	return q
}

func (q *Query) Where(where map[string]interface{}) *Query {
	q.where = where
	return q
}

func (q *Query) Group(fields...string) *Query {
	q.group = "GROUP BY "+strings.Join(fields, ",")
	return q
}

func (q *Query) Having(having string) *Query {
	q.having = "HAVING "+ having
	return q
}

func (q *Query) Order(orders...string) *Query {
	q.order = "ORDER BY "+ strings.Join(orders, ",")
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

func BuildQuery(q *Query) (sql string, args []interface{}) {
	args  = utils.AcquireArgs()
	where := ""

	condition := buildWhere(q.where, args)
	if len(condition) > 0 {
		where = string(condition)
	}

	return fmt.Sprintf(selectSql, q.columns, q.table, where, q.group, q.having, q.order, q.offset, q.limit), args
}
