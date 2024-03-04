package argument

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

const (
	OperatorTypeEqual    = "equal"
	OperatorTypeNotEqual = "not_equal"
	OperatorTypeGT       = "gt"  // greater than, >
	OperatorTypeGTE      = "gte" // greater than or equal, >=
	OperatorTypeLT       = "lt"  // less than, <
	OperatorTypeLTE      = "lte" // less than or equal, <=
	OperatorTypeIn       = "in"
	OperatorTypeNotIn    = "not_in"
)

type Operation interface {
	ToSql() string
	validateAndFormat() error
}

func OperationFactory(operatorType string, fieldName string, value interface{}) (Operation, error) {
	switch operatorType {
	case OperatorTypeEqual, OperatorTypeNotEqual, OperatorTypeGT, OperatorTypeGTE, OperatorTypeLT, OperatorTypeLTE:
		return newCompareOperation(fieldName, operatorType, value)
	case OperatorTypeIn, OperatorTypeNotIn:
		return newContainsOperation(fieldName, operatorType, value)
	default:
		return nil, fmt.Errorf("unsupported operation type: %s", operatorType)
	}
}

var _ Operation = (*CompareOperation)(nil)
var _ Operation = (*ContainsOperation)(nil)

// CompareOperation represents an operation in the SQL statement for comparison.
type CompareOperation struct {
	fieldName string
	operator  string
	value     interface{}
}

func newCompareOperation(fieldName string, operator string, value interface{}) (Operation, error) {
	op := &CompareOperation{
		fieldName: fieldName,
		operator:  operator,
		value:     value,
	}
	if err := op.validateAndFormat(); err != nil {
		return nil, err
	}
	return op, nil
}

func (e *CompareOperation) ToSql() string {
	return fmt.Sprintf(" %s %s '%v' ", e.fieldName, e.getOperator(), e.value)
}

func (e *CompareOperation) validateAndFormat() error {
	if e.value == nil {
		return fmt.Errorf("CompareOperation.value cannot be nil")
	}
	// Based on the function parseAstValue, most data types are converted to String,
	// except for the Bool type.
	// Therefore, the validateAndFormat logic here should ensure that the value's data type is either String or Bool.
	// Other types should result in an error directly.
	value := reflect.ValueOf(e.value)
	valueType := value.Type()

	switch valueType.Kind() {
	case reflect.String:
		return nil
	case reflect.Bool:
		if value.Bool() {
			e.value = 1
		} else {
			e.value = 0
		}
		return nil
	default:
		return fmt.Errorf("CompareOperation expects the value to be a string or bool, but got %s ", valueType.String())
	}
}

func (e *CompareOperation) getOperator() string {
	switch e.operator {
	case OperatorTypeEqual:
		return "="
	case OperatorTypeNotEqual:
		return "!="
	case OperatorTypeGT:
		return ">"
	case OperatorTypeGTE:
		return ">="
	case OperatorTypeLT:
		return "<"
	case OperatorTypeLTE:
		return "<="
	default:
		slog.Warn(fmt.Sprintf("unsupported operator type: %s", e.operator))
		return ""
	}
}

type ContainsOperation struct {
	fieldName   string
	operator    string
	value       interface{}
	innerValues []interface{}
}

func newContainsOperation(fieldName string, operator string, value interface{}) (*ContainsOperation, error) {
	op := &ContainsOperation{
		fieldName: fieldName,
		operator:  operator,
		value:     value,
	}
	if err := op.validateAndFormat(); err != nil {
		return nil, err
	}
	return op, nil
}

func (c *ContainsOperation) validateAndFormat() error {
	if c.value == nil {
		return fmt.Errorf("ContainsOperation.value cannot be nil")
	}

	value := reflect.ValueOf(c.value)
	valueType := value.Type()

	if valueType.Kind() == reflect.Slice || valueType.Kind() == reflect.Array {
		c.innerValues = make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			c.innerValues[i] = value.Index(i).Interface()
		}
		return nil
	}
	return fmt.Errorf("ContainsOperation expects the value to be an array or slice, but got %s ", valueType.String())
}

func (c *ContainsOperation) ToSql() string {
	if c.innerValues == nil || len(c.innerValues) == 0 {
		return ""
	}
	values := make([]string, len(c.innerValues))
	for i, v := range c.innerValues {
		values[i] = fmt.Sprintf("'%v'", v)
	}

	return fmt.Sprintf(" %s %s (%s) ", c.fieldName, c.getOperator(), strings.Join(values, ","))
}

func (c *ContainsOperation) getOperator() string {
	switch c.operator {
	case OperatorTypeIn:
		return "IN"
	case OperatorTypeNotIn:
		return "NOT IN"
	default:
		slog.Warn(fmt.Sprintf("unsupported operator type: %s", c.operator))
		return ""
	}

}
