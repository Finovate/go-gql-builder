package model

import (
	"github.com/graphql-go/graphql"

	"github.com/Finovate/go-gql-builder/pkg/adapter"
	"github.com/Finovate/go-gql-builder/pkg/core"
	"github.com/Finovate/go-gql-builder/pkg/core/argument"
)

const (
	FieldTypeUser = "user"
)

type UserDelegate struct {
	adapter.SqlAdapter
	core.BaseNode
	argument.DefaultArgumentBuilder
}

var _ core.Node = (*UserDelegate)(nil)

func NewUserDelegate() (d *UserDelegate) {
	d = &UserDelegate{}
	d.SqlAdapter = adapter.NewDefaultSqlAdapter("user", d.initItemTable(), d)
	return
}

func (d *UserDelegate) initItemTable() []*adapter.Column {
	columns := make([]*adapter.Column, 0)
	column := &adapter.Column{
		Type:  "",
		Name:  "id",
		Alias: "id",
	}
	column.SetPrimaryKey()
	columns = append(columns, column)

	columns = append(columns, &adapter.Column{
		Type:  "",
		Name:  "name",
		Alias: "name",
	})
	columns = append(columns, &adapter.Column{
		Type:  "",
		Name:  "email",
		Alias: "email",
	})

	return columns
}

func (d *UserDelegate) Name() string {
	return "users"
}

func (d *UserDelegate) Type() core.FieldType {
	return FieldTypeUser
}

func (d *UserDelegate) IsList() bool {
	return true
}

func (d *UserDelegate) BuildFields() []*core.Field {
	fields := make([]*core.Field, 0)

	fields = append(fields,
		core.NewNodeField("id", core.FieldTypeString),
	)
	fields = append(fields,
		core.NewNodeField("name", core.FieldTypeString),
	)
	fields = append(fields,
		core.NewNodeField("price", core.FieldTypeFloat),
	)

	fields = append(fields, d.departmentField())
	return fields
}

func (d *UserDelegate) departmentField() *core.Field {
	field := core.NewNodeField("department", FieldTypeDepartment)
	field.SetResolver(func(p graphql.ResolveParams) (interface{}, error) {
		// TODO 可以考虑封装一个Thunk，这个resolver只拿主键，然后外层的方法再去查具体数据
		// github.com/graphql-go/graphql@v0.8.1/executor.go:754
		return []map[string]interface{}{{"id": "1", "name": "Example Product", "price": 99.99}}, nil
	})
	return field
}
