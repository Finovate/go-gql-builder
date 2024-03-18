package argument

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core/argument"
)

var (
	filterArgumentType = graphql.NewScalar(graphql.ScalarConfig{
		Name:        FilterArgumentType,
		Description: "Filter argument",
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

var _ SqlArgument = (*FilterArgument)(nil)

type FilterArgument struct {
	operationsMap map[string][]Operation
}

func newFilterArgument() argument.Argument {
	return &FilterArgument{
		operationsMap: make(map[string][]Operation),
	}
}

func (f *FilterArgument) TypeName() string {
	return FilterArgumentType
}

func (f *FilterArgument) Validate(input interface{}) error {
	argsMap, ok := input.(map[string]interface{})
	if !ok {
		return fmt.Errorf("filter argument must be a map[string]interface{}")
	}
	for fieldName, rawMap := range argsMap {
		operationMap, ok := rawMap.(map[string]interface{})
		if !ok {
			return fmt.Errorf("argument for field %s must be a map[string]interface{}", fieldName)
		}

		for op, value := range operationMap {
			operation, err := OperationFactory(op, fieldName, value)
			if err != nil {
				return err
			}
			if f.operationsMap[fieldName] == nil {
				f.operationsMap[fieldName] = make([]Operation, 0)
			}
			f.operationsMap[fieldName] = append(f.operationsMap[fieldName], operation)
		}

	}

	// TODO validateAndFormat fieldName in Node fields or table columns
	return nil
}

func (f *FilterArgument) GetArgumentType() graphql.Input {
	return filterArgumentType
}

func (f *FilterArgument) ParseSqlValue() string {
	sqlStrings := make([]string, 0, len(f.operationsMap))
	for _, opList := range f.operationsMap {
		for _, operation := range opList {
			sqlStrings = append(sqlStrings, operation.ToSql())
		}
	}

	return strings.Join(sqlStrings, " AND ")
}

func (f *FilterArgument) CombineSql(clauses *QueryClauses) {
	clauses.SetWhere(f.ParseSqlValue())
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
		return valueAST.Value
	case *ast.FloatValue:
		return valueAST.Value
	case *ast.EnumValue:
		return valueAST.Value
	default:
		return nil
	}
}

func init() {
	argument.RegisterArgument(FilterArgumentType, newFilterArgument)
}
