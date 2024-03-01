package core

import (
	"database/sql"
	"fmt"

	"github.com/graphql-go/graphql"
)

// NodeRegistry is a key component in the go-gql-builder framework,
// serving as a repository for storing and managing the implementations of Node declared by developers.
type NodeRegistry struct {
	nodes       []Node
	nodesByType map[FieldType]Node

	fieldsMap map[FieldType]graphql.Fields
	argsMap   map[FieldType]graphql.FieldConfigArgument

	// 用一个缓存先初始化所有的node, 以免在具体构建field时依赖了一个不存在的node.
	// 比如user.department 依赖了 department 这个node，在处理这个field的时候如果没有预创建这一步，
	// 那么就会导致找不到对应的类型，从而导致field创建失败.
	// 在最终处理node的时候，也应该从缓存中加载相应的指针出来进行最终的构建.
	preCache map[FieldType]graphql.Output

	completeCache graphql.Fields

	// TODO  HubSet 框架支持多个数据源
	db *sql.DB
}

func NewRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes:         make([]Node, 0),
		nodesByType:   make(map[FieldType]Node),
		fieldsMap:     make(map[FieldType]graphql.Fields),
		argsMap:       make(map[FieldType]graphql.FieldConfigArgument),
		preCache:      make(map[FieldType]graphql.Output),
		completeCache: make(graphql.Fields),
	}
}

func DefaultRegistry() *NodeRegistry {
	if registry == nil {
		registry = NewRegistry()
	}
	return registry
}

func (h *NodeRegistry) GetDB() *sql.DB {
	return h.db
}

func (h *NodeRegistry) SetDB(db *sql.DB) {
	h.db = db
}

func (h *NodeRegistry) Register(delegate Node) {
	h.nodes = append(h.nodes, delegate)
	h.nodesByType[delegate.Type()] = delegate
	delegate.SetRegistry(h)
}

func (h *NodeRegistry) getNode(typeName FieldType) (Node, error) {
	node, ok := h.nodesByType[typeName]
	if !ok {
		return nil, fmt.Errorf("unsupported node type: %s", typeName)
	}
	return node, nil
}

func (h *NodeRegistry) BuildSchema() (*graphql.Schema, error) {

	h.preLoadDelegate()

	for _, delegate := range h.nodes {
		err := h.buildNode(delegate)
		if err != nil {
			return nil, err
		}
	}

	// 生成schema(逻辑不变)
	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: h.completeCache,
		},
	)

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

func (h *NodeRegistry) preLoadDelegate() {
	// 预加载delegate
	for _, delegate := range h.nodes {
		obj := graphql.NewObject(graphql.ObjectConfig{
			Name:   delegate.Name(),
			Fields: make(graphql.Fields),
		})

		var result graphql.Output
		result = obj
		if delegate.IsList() {
			result = graphql.NewList(obj)
		}

		h.preCache[delegate.Type()] = result
	}
}

// initNodeField Node由一组Field组成，这个方法中会解析Node下面的Field，并将其转换为graphql.Field以及Argument
func (h *NodeRegistry) initNodeField(delegate Node) error {
	rawFields := delegate.BuildFields()
	fields := make(graphql.Fields)
	args := make(graphql.FieldConfigArgument)
	for _, f := range rawFields {
		convert, err := f.Convert(h)
		if err != nil {
			return err
		}
		fields[f.fieldName] = convert
	}
	argList := delegate.BuildArgs()
	for _, arg := range argList {
		args[arg.TypeName()] = &graphql.ArgumentConfig{Type: arg.GetArgumentType()}
	}

	h.fieldsMap[delegate.Type()] = fields
	h.argsMap[delegate.Type()] = args
	return nil
}

// buildNode 入参是一个Node，该方法将会根据Node的信息对应的schema(graphql.Field)同时存进completeCache中
// 如:
//
//	users{
//	  id
//	  name
//	  department{
//	    id
//	    name
//	  }
//	}
//
// 同时，这个schema也需要能够根据args进行查询过滤，即users(id: "123")
// 这个方法应该是一个递归的过程，users中的一个field依赖了Node department，那么build的顺序应该是users -> users.department -> department
func (h *NodeRegistry) buildNode(delegate Node) error {
	if _, ok := h.completeCache[delegate.Name()]; ok {
		return nil
	}
	err := h.initNodeField(delegate)
	if err != nil {
		return err
	}

	var obj *graphql.Object

	cache := h.preCache[delegate.Type()]
	if tmp, ok := cache.(*graphql.List); ok {
		obj = tmp.OfType.(*graphql.Object)
	} else if tmp, ok := cache.(*graphql.Object); ok {
		obj = tmp
	}

	if obj == nil {
		return fmt.Errorf("unsupported field type: %s", delegate.Type())
	}

	fields := h.fieldsMap[delegate.Type()]
	args := h.argsMap[delegate.Type()]
	for name, field := range fields {
		obj.AddFieldConfig(name, field)
	}

	h.completeCache[delegate.Name()] = &graphql.Field{
		Type:    cache,
		Args:    args,
		Resolve: delegate.Resolve(),
	}
	return nil
}

func (h *NodeRegistry) loadFieldType(flag FieldType) (out graphql.Output, isDefaultFieldType bool, err error) {
	if fieldType, ok := defaultFieldTypeMapping[flag]; ok {
		return fieldType, true, nil
	}

	// 自定义的fieldType，尝试从cache中加载
	if fieldType, ok := h.preCache[flag]; ok {
		return fieldType, false, nil
	}

	return nil, false, fmt.Errorf("unsupported field type: %s", flag)
}

var registry *NodeRegistry

var defaultFieldTypeMapping = map[FieldType]graphql.Output{
	FieldTypeString:  graphql.String,
	FieldTypeInt:     graphql.Int,
	FieldTypeFloat:   graphql.Float,
	FieldTypeBoolean: graphql.Boolean,
}
