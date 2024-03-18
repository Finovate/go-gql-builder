package argument

import (
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

type Builder func() Argument

var argumentBuilders = make(map[string]Builder)

func RegisterArgument(typename string, arg Builder) {
	argumentBuilders[typename] = arg
}

type DefaultArgumentBuilder struct {
}

func (i *DefaultArgumentBuilder) BuildArgs() []Argument {
	res := make([]Argument, 0, len(argumentBuilders))
	for _, builder := range argumentBuilders {
		res = append(res, builder())
	}
	return res
}

// Factory create a argument instance by typename, when typename is not exist, return nil.
func Factory(typename string) Argument {
	builder, ok := argumentBuilders[typename]
	if ok {
		return builder()
	}
	return nil
}
