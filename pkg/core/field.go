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
// 初始化结束后从hub的argsMap中获取这个Node的ArgumentConfig，并将其作为Field的Args属性。
// 最终的结果将会是嵌套的field也能够支持Argument参数
func (f *Field) Convert(hub *NodeRegistry) (field *graphql.Field, argConfig graphql.FieldConfigArgument, err error) {
	// 一个field有可能回依赖其他node对象
	t, isDefault, err := hub.loadFieldType(f.fieldType)
	if err != nil {
		return nil, nil, err
	}

	field = &graphql.Field{
		Name:    f.fieldName,
		Type:    t,
		Resolve: f.resolver,
	}

	// 当field的类型是默认类型时，生成ArgumentConfig
	if isDefault {
		argConfig = graphql.FieldConfigArgument{
			f.fieldName: &graphql.ArgumentConfig{
				Type: t,
			},
		}
		return field, argConfig, nil
	}

	// 当field是一个Node自定义的fieldType的时候，是不是应该递归的初始化某种东西？
	node, err := hub.getNode(f.fieldType)
	if err != nil {
		return nil, nil, err
	}

	if err = hub.buildNode(node); err != nil {
		return nil, nil, err
	}
	argConfig = hub.argsMap[node.Type()]
	field.Args = argConfig
	return field, nil, nil
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
