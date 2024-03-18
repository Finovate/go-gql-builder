package argument

import (
	"fmt"

	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core/argument"
)

type SqlArgument interface {
	argument.Argument
	ParseSqlValue() string
	CombineSql(clauses *QueryClauses)
}

type QueryClauses struct {
	selectColumn string
	from         string
	where        string
	groupBy      string
	orderBy      string
	limit        string
}

func NewQueryClauses(columns, db string) *QueryClauses {
	return &QueryClauses{selectColumn: columns, from: db}
}

func (c *QueryClauses) SetSelect(columns string) {
	c.selectColumn = columns
}

func (c *QueryClauses) SetFrom(db string) {
	c.from = db
}

func (c *QueryClauses) SetWhere(filter string) {
	c.where = filter
}

func (c *QueryClauses) SetGroupBy(g string) {
	c.groupBy = g
}

func (c *QueryClauses) SetOrderBy(o string) {
	c.orderBy = o
}

func (c *QueryClauses) SetLimit(l string) {
	c.limit = l
}

func (c *QueryClauses) ToSql() (string, error) {
	sql := ""
	if c.selectColumn == "" || c.from == "" {
		return "", fmt.Errorf("not enough fields combined to form SQL statements")
	}
	sql += fmt.Sprintf("SELECT %s FROM %s", c.selectColumn, c.from)

	if c.where != "" {
		sql += fmt.Sprintf(" WHERE %s", c.where)
	}
	if c.groupBy != "" {
		sql += fmt.Sprintf(" GroupBy %s", c.groupBy)
	}
	if c.orderBy != "" {
		sql += fmt.Sprintf(" OrderBy %s", c.orderBy)
	}
	if c.limit != "" {
		sql += fmt.Sprintf(" LIMIT %s", c.limit)
	}

	return sql, nil
}
