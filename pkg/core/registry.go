package core

import (
	"database/sql"
	"fmt"
	"github.com/graphql-go/graphql"
)

// NodeRegistry is a key component in the go-gql-builder framework,
// serving as a repository for storing and managing the implementations of Node declared by developers.
type NodeRegistry struct {
	nodes []Node

	// 用一个缓存先初始化所有的node, 以免在具体构建field时依赖了一个不存在的node.
	// 比如user.department 依赖了 department 这个node，在处理这个field的时候如果没有预创建这一步，
	// 那么就会导致找不到对应的类型，从而导致field创建失败.
	// 在最终处理node的时候，也应该从缓存中加载相应的指针出来进行最终的构建.
	preCache map[FieldType]graphql.Output

	// TODO  HubSet 框架支持多个数据源
	db *sql.DB
}

func Registry() *NodeRegistry {
	if registry == nil {
		registry = new(NodeRegistry)
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
}

func (h *NodeRegistry) BuildSchema() (*graphql.Schema, error) {

	h.preLoadDelegate()

	fields := make(graphql.Fields)
	for _, delegate := range h.nodes {
		node, err := h.buildNode(delegate)
		if err != nil {
			return nil, err
		}

		fields[delegate.Name()] = node
	}

	// 生成schema(逻辑不变)
	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: fields,
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
	h.preCache = make(map[FieldType]graphql.Output)
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

// 一个field依赖其他delegate对象，就发生在这里
func (h *NodeRegistry) initNodeField(delegate Node) (graphql.Fields, error) {
	rawFields := delegate.BuildField()
	fields := make(graphql.Fields)
	for _, f := range rawFields {
		convert, err := f.Convert(h)
		if err != nil {
			return nil, err
		}
		fields[f.fieldName] = convert
	}
	return fields, nil
}

// 这个方法最终应该输出一个能够被Hub直接使用的field字段
func (h *NodeRegistry) buildNode(delegate Node) (*graphql.Field, error) {
	fields, err := h.initNodeField(delegate)
	if err != nil {
		return nil, err
	}

	var obj *graphql.Object

	cache := h.preCache[delegate.Type()]
	if tmp, ok := cache.(*graphql.List); ok {
		obj = tmp.OfType.(*graphql.Object)
	} else if tmp, ok := cache.(*graphql.Object); ok {
		obj = tmp
	}

	if obj == nil {
		return nil, fmt.Errorf("unsupported field type: %s", delegate.Type())
	}

	for name, field := range fields {
		obj.AddFieldConfig(name, field)
	}

	return &graphql.Field{
		Type:    cache,
		Resolve: delegate.Resolve(),
	}, nil
}

func (h *NodeRegistry) loadFieldType(flag FieldType) (graphql.Output, error) {
	if fieldType, ok := defaultFieldTypeMapping[flag]; ok {
		return fieldType, nil
	}

	// 自定义的fieldType，尝试从cache中加载
	if fieldType, ok := h.preCache[flag]; ok {
		return fieldType, nil
	}

	return nil, fmt.Errorf("unsupported field type: %s", flag)
}

var registry *NodeRegistry

var defaultFieldTypeMapping = map[FieldType]graphql.Output{
	FieldTypeString:  graphql.String,
	FieldTypeInt:     graphql.Int,
	FieldTypeFloat:   graphql.Float,
	FieldTypeBoolean: graphql.Boolean,
}
