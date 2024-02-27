package argument

import "fmt"

const (
	OperationTypeEqual = "equal"
)

type Operation interface {
	ToSql() string
	Validate() error
}

func OperationFactory(operationType string, fieldName string, value interface{}) (Operation, error) {
	switch operationType {
	case OperationTypeEqual:
		return &EqualOperation{fieldName, value}, nil
		//TODO add more operation types
	default:
		return nil, fmt.Errorf("unsupported operation type: %s", operationType)
	}
}

var _ Operation = (*EqualOperation)(nil)

type EqualOperation struct {
	fieldName string
	Value     interface{}
}

func (e *EqualOperation) ToSql() string {
	return fmt.Sprintf(" %s = '%s' ", e.fieldName, e.Value)
}

func (e *EqualOperation) Validate() error {
	//
	return nil
	//TODO implement me
	panic("implement me")
}
