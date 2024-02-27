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

type SqlArgument interface {
	Argument
	ParseSqlValue() string
}

type DefaultArgumentBuilder struct {
}

func (i *DefaultArgumentBuilder) BuildArgs() []Argument {
	return []Argument{newFilterArgument()}
}

// ArgumentFactory create a argument instance by typename, when typename is not exist, return nil.
func ArgumentFactory(typename string) Argument {
	switch typename {
	case FilterArgumentType:
		return newFilterArgument()
	default:
		return nil
	}
}
