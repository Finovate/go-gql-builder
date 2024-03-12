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
}

type QueryClauses struct {
	SELECT  string
	FROM    string
	WHERE   string
	GroupBy string
	OrderBy string
	LIMIT   string
	FOR     string
}

func (c *QueryClauses) ToSql() string {
	sql := ""
	if c.SELECT != "" {
		sql += fmt.Sprintf("SELECT %s", c.SELECT)
	}
	if c.FROM != "" {
		sql += fmt.Sprintf(" FROM %s", c.FROM)
	}
	if c.WHERE != "" {
		sql += fmt.Sprintf(" WHERE %s", c.WHERE)
	}
	if c.GroupBy != "" {
		sql += fmt.Sprintf(" GroupBy %s", c.GroupBy)
	}
	if c.GroupBy != "" {
		sql += fmt.Sprintf(" GroupBy %s", c.GroupBy)
	}
	if c.LIMIT != "" {
		sql += fmt.Sprintf(" LIMIT %s", c.LIMIT)
	}
	if c.FOR != "" {
		sql += fmt.Sprintf(" FOR %s", c.FOR)
	}

	return sql
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
