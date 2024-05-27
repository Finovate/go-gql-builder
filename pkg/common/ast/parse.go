package ast

import (
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql/language/ast"
)

// ParseAstValue 辅助函数：递归解析 AST 值
// TODO: 基本数据类型的解析逻辑还有点粗糙，需要优化。这个函数优化后还需要修改 CompareOperation
func ParseAstValue(valueAST ast.Value) interface{} {
	switch valueAST := valueAST.(type) {
	case *ast.ObjectValue:
		result := make(map[string]interface{})
		for _, field := range valueAST.Fields {
			result[field.Name.Value] = ParseAstValue(field.Value)
		}
		return result
	case *ast.ListValue:
		list := make([]interface{}, len(valueAST.Values))
		for i, value := range valueAST.Values {
			list[i] = ParseAstValue(value)
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
