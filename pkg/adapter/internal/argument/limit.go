package argument

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core/argument"
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
			return parseAstValue(valueAST)
		},
	})
)

var _ SqlArgument = (*LimitArgument)(nil)

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
	err := fmt.Errorf("limit argument must be a positive integer or an all positive integer array with a length of 2")

	switch arg := input.(type) {
	case int:
		f.limit = arg
	case []interface{}:
		if len(arg) != 2 {
			return err
		}
		for _, i := range arg {
			if _, ok := i.(int); !ok {
				return err
			}
		}
		f.offset = arg[0].(int)
		f.limit = arg[1].(int)
	default:
		return err
	}

	if f.limit >= 0 && f.offset >= 0 {
		return nil
	}
	return err
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
