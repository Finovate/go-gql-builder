package argument

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	astCommon "github.com/Finovate/go-gql-builder/pkg/common/ast"
	"github.com/Finovate/go-gql-builder/pkg/core/argument"
)

var (
	orderByArgumentType = graphql.NewScalar(graphql.ScalarConfig{
		Name:        OrderByArgumentType,
		Description: "OrderBy argument",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			return astCommon.ParseAstValue(valueAST)
		},
	})
)

var _ SqlArgument = (*OrderByArgument)(nil)

var (
	sort = map[string]struct{}{
		"ASC":  {},
		"DESC": {},
	}
)

type OrderByArgument struct {
	sortMap map[string]string
}

func newOrderByArgument() argument.Argument {
	return &OrderByArgument{
		sortMap: make(map[string]string, 0),
	}
}

func (f *OrderByArgument) TypeName() string {
	return OrderByArgumentType
}

func (f *OrderByArgument) Validate(input interface{}) error {
	argsMap, ok := input.(map[string]interface{})
	if !ok {
		return fmt.Errorf("orderBy argument must be a map[string]string{}")
	}

	for fieldName, rawString := range argsMap {
		sortString, ok := rawString.(string)
		if !ok {
			return fmt.Errorf("argument for field %s must be a string", fieldName)
		}

		_, ok = sort[strings.ToUpper(sortString)]
		if !ok {
			return fmt.Errorf(`argument for field %s must be "asc" or "desc"`, fieldName)
		}
		f.sortMap[fieldName] = sortString
	}

	return nil
}

func (f *OrderByArgument) GetArgumentType() graphql.Input {
	return orderByArgumentType
}

func (f *OrderByArgument) ParseSqlValue() string {
	sqlStrings := make([]string, 0, len(f.sortMap))
	for fieldName, sort := range f.sortMap {
		sqlStrings = append(sqlStrings, fmt.Sprintf("%s %s", fieldName, sort))
	}

	return strings.Join(sqlStrings, ",")
}

func (f *OrderByArgument) CombineSql(clauses *QueryClauses) {
	clauses.SetOrderBy(f.ParseSqlValue())
}

func init() {
	argument.RegisterArgument(OrderByArgumentType, newOrderByArgument)
}
