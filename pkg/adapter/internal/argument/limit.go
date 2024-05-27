package argument

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	astCommon "github.com/Finovate/go-gql-builder/pkg/common/ast"
	"github.com/Finovate/go-gql-builder/pkg/core/argument"
)

var (
	limitArgumentType = graphql.NewScalar(graphql.ScalarConfig{
		Name:        LimitArgumentType,
		Description: "Limit argument",
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

var _ SqlArgument = (*LimitArgument)(nil)

// LimitArgument
// limit 入参
// limit:{ count: 10 offset : 0 }
type LimitArgument struct {
	limit  int
	offset int
}

func newLimitArgument() argument.Argument {
	return &LimitArgument{
		limit:  0,
		offset: 0,
	}
}

func (f *LimitArgument) TypeName() string {
	return LimitArgumentType
}

func (f *LimitArgument) Validate(input interface{}) error {
	argsMap, ok := input.(map[string]interface{})
	if !ok {
		return fmt.Errorf("limit argument must be a map[string]int")
	}

	f.offset, _ = argsMap["offset"].(int)
	if f.limit, ok = argsMap["count"].(int); !ok {
		return fmt.Errorf("limit argument field count must be required")
	}

	if f.limit >= 0 && f.offset >= 0 {
		return nil
	} else {
		return fmt.Errorf("limit argument must be a positive integer")
	}
}

func (f *LimitArgument) GetArgumentType() graphql.Input {
	return limitArgumentType
}

func (f *LimitArgument) ParseSqlValue() string {
	return fmt.Sprintf("%v,%v", f.offset, f.limit)
}

func (f *LimitArgument) CombineSql(clauses *QueryClauses) {
	clauses.SetLimit(f.ParseSqlValue())
}

func init() {
	argument.RegisterArgument(LimitArgumentType, newLimitArgument)
}
