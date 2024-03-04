package core

import (
	"github.com/graphql-go/graphql"
)

// Field a Field represents a component of the Node abstraction,
// corresponding to an attribute of a business object.
// Field is used to define and represent a single data point within a Node,
// such as a user's name, age, etc.
type Field struct {
	fieldName string
	fieldType FieldType

	asList bool

	resolver graphql.FieldResolveFn
}

// Convert Field对象不仅需要转换成graphql.Field对象，同时要根据自身的数据类型，生成相应的ArgumentConfig
// 当 Field 是一个对象类型时，将会递归调用 NodeRegistry.buildNode 方法，优先初始化对应的Node，
func (f *Field) Convert(hub *NodeRegistry) (field *graphql.Field, err error) {
	// 一个field有可能回依赖其他node对象
	t, isDefault, err := hub.loadFieldType(f.fieldType)
	if err != nil {
		return nil, err
	}

	field = &graphql.Field{
		Name:    f.fieldName,
		Type:    t,
		Resolve: f.resolver,
	}

	// 当field的类型是默认类型时
	if !isDefault {
		// When the field type is a custom Node type, recursively initialize the Node.
		node, err := hub.getNode(f.fieldType)
		if err != nil {
			return nil, err
		}

		if err = hub.buildNode(node); err != nil {
			return nil, err
		}

		field.Args = hub.argsMap[node.Type()]
	}

	return field, nil
}

func (f *Field) SetResolver(resolver graphql.FieldResolveFn) {
	f.resolver = resolver
}

func (f *Field) Resolver() graphql.FieldResolveFn {
	return f.resolver
}

func NewNodeField(fieldName string, fieldType FieldType) *Field {
	return &Field{
		fieldName: fieldName,
		fieldType: fieldType,
		resolver:  nil,
	}
}
