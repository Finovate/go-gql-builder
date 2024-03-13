package argument

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

// Argument represents an abstracted definition of GraphQL's conditional querying capabilities.
// It is decoupled from specific Query implementations,
// offering a universal and flexible solution for conditional queries.
// This approach transforms Argument into a reusable component,
// standardizing condition filtering and data processing across various querying contexts,
// thereby enhancing the framework's versatility and adaptability.
type Argument interface {
	TypeName() string
	// Validate checks if the value that pass in by client is legal.
	Validate(interface{}) error
	GetArgumentType() graphql.Input
}

type SqlArgument interface {
	Argument
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

type DefaultArgumentBuilder struct {
}

func (i *DefaultArgumentBuilder) BuildArgs() []Argument {
	return []Argument{newFilterArgument()}
}

// ArgumentFactory create a argument instance by typename, when typename is not exist, return nil.
func ArgumentFactory(typename string) Argument {
	switch typename {
	case FilterArgumentType:
		return newFilterArgument()
	default:
		return nil
	}
}
