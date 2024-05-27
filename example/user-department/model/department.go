package model

import (
	"github.com/Finovate/go-gql-builder/pkg/adapter"
	"github.com/Finovate/go-gql-builder/pkg/core"
	"github.com/Finovate/go-gql-builder/pkg/core/argument"
)

const (
	FieldTypeDepartment = "department"
)

type DepartmentDelegate struct {
	adapter.SqlAdapter
	core.BaseNode
	argument.DefaultArgumentBuilder
}

var _ core.Node = (*DepartmentDelegate)(nil)

func (d *DepartmentDelegate) Name() string {
	return "departments"
}

func (d *DepartmentDelegate) Type() core.FieldType {
	return FieldTypeDepartment
}

func (d *DepartmentDelegate) BuildFields() []*core.Field {
	fields := make([]*core.Field, 0)

	fields = append(fields,
		core.NewNodeField("id", core.FieldTypeString),
	)
	fields = append(fields,
		core.NewNodeField("name", core.FieldTypeString),
	)
	fields = append(fields,
		core.NewNodeField("test", core.FieldTypeString),
	)
	return fields
}

func (d *DepartmentDelegate) IsList() bool {
	return true
}

func NewDepartmentDelegate() (d *DepartmentDelegate) {
	d = &DepartmentDelegate{}
	d.SqlAdapter = adapter.NewDefaultSqlAdapter("department", d.initItemTable(), d)
	return
}

func (d *DepartmentDelegate) initItemTable() []*adapter.Column {
	columns := make([]*adapter.Column, 0)
	columns = append(columns, &adapter.Column{
		Type:  "",
		Name:  "id",
		Alias: "id",
	})
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
