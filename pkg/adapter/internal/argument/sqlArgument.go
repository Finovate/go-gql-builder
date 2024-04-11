package argument

import (
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql/language/ast"

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
		sql += fmt.Sprintf(" Group By %s", c.groupBy)
	}
	if c.orderBy != "" {
		sql += fmt.Sprintf(" Order By %s", c.orderBy)
	}
	if c.limit != "" {
		sql += fmt.Sprintf(" LIMIT %s", c.limit)
	}

	return sql, nil
}

// 辅助函数：递归解析 AST 值
// TODO: 基本数据类型的解析逻辑还有点粗糙，需要优化。这个函数优化后还需要修改 CompareOperation
func parseAstValue(valueAST ast.Value) interface{} {
	switch valueAST := valueAST.(type) {
	case *ast.ObjectValue:
		result := make(map[string]interface{})
		for _, field := range valueAST.Fields {
			result[field.Name.Value] = parseAstValue(field.Value)
		}
		return result
	case *ast.ListValue:
		list := make([]interface{}, len(valueAST.Values))
		for i, value := range valueAST.Values {
			list[i] = parseAstValue(value)
		}
		return list
	case *ast.StringValue:
		return valueAST.Value
	case *ast.BooleanValue:
		return valueAST.Value
	case *ast.IntValue:
		result, err := strconv.Atoi(valueAST.Value)
		if err != nil {
			fmt.Println(err)
			return valueAST.Value
		}
		return result
	case *ast.FloatValue:
		return valueAST.Value
	case *ast.EnumValue:
		return valueAST.Value
	default:
		return nil
	}
}
