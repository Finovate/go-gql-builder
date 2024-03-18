package argument

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsOperation_ValidateAndToSql(t *testing.T) {
	// normal cases
	executionCases := []ContainsOperation{
		{
			fieldName: "age",
			operator:  OperatorTypeIn,
			value:     []int{1, 2, 3},
		},
		{
			fieldName: "age",
			operator:  OperatorTypeIn,
			value:     []string{"a", "b", "c"},
		},
		{
			fieldName: "age",
			operator:  OperatorTypeNotIn,
			value:     []float64{1.1, 2.2, 3.3},
		},
	}

	exceptionSql := []string{
		" age IN ('1','2','3') ",
		" age IN ('a','b','c') ",
		" age NOT IN ('1.1','2.2','3.3') ",
	}

	for i := 0; i < len(executionCases); i++ {
		op := executionCases[i]

		curOp, err := OperationFactory(op.operator, op.fieldName, op.value)
		if err != nil {
			t.Fatalf("Validate failed for case %d: %v", i, err)
			return
		}
		actualSql := curOp.ToSql()
		assert.Equalf(t, exceptionSql[i], actualSql, "ToSql() failed for case %d", i)
	}

	// error case
	errorCases := []ContainsOperation{
		{
			fieldName: "age",
			operator:  OperatorTypeIn,
			value:     nil,
		},
		{
			fieldName: "age",
			operator:  OperatorTypeIn,
			value:     "abc",
		},
	}
	for i := 0; i < len(errorCases); i++ {
		op := errorCases[i]
		_, err := OperationFactory(op.operator, op.fieldName, op.value)
		if err == nil {
			t.Fatalf("Validate should fail for case %d", i)
		}
	}
}

func TestCompareOperation_ValidateAndSql(t *testing.T) {
	// normal cases
	executionCases := []ContainsOperation{
		{
			fieldName: "age",
			operator:  OperatorTypeEqual,
			value:     "18",
		},
		{
			fieldName: "is_male",
			operator:  OperatorTypeNotEqual,
			value:     true,
		},
	}

	exceptionSql := []string{
		" age = '18' ",
		" is_male != '1' ",
	}

	for i := 0; i < len(executionCases); i++ {
		op := executionCases[i]

		curOp, err := OperationFactory(op.operator, op.fieldName, op.value)
		if err != nil {
			t.Fatalf("Validate failed for case %d: %v", i, err)
			return
		}
		actualSql := curOp.ToSql()
		assert.Equalf(t, exceptionSql[i], actualSql, "ToSql() failed for case %d", i)
	}

	// error case
	errorCases := []ContainsOperation{
		{
			fieldName: "age",
			operator:  OperatorTypeEqual,
			value:     nil,
		},
		{
			fieldName: "age",
			operator:  OperatorTypeEqual,
			value:     []string{"abc"},
		},
	}
	for i := 0; i < len(errorCases); i++ {
		op := errorCases[i]
		_, err := OperationFactory(op.operator, op.fieldName, op.value)
		if err == nil {
			t.Fatalf("Validate should fail for case %d", i)
		}
	}

}
