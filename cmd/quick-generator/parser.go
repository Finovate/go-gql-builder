package main

import (
	"bytes"
	"fmt"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"
	"strings"
	"text/template"
)

var templateString = `package model

import (
	"github.com/graphql-go/graphql"

	"github.com/shuishiyuanzhong/go-gql-builder/pkg/adapter"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core/argument"
)

const (
	FieldType{{ .NodeName }} = "{{ .NodeNameLower }}"
)

type {{ .NodeName }} struct {
	adapter.SqlAdapter
	core.BaseNode
	argument.DefaultArgumentBuilder
}

var _ core.Node = (*{{ .NodeName }})(nil)

func New{{ .NodeName }}() (d *{{ .NodeName }}) {
	d = &{{ .NodeName }}{}
	d.SqlAdapter = adapter.NewDefaultSqlAdapter("{{ .TableName }}", d.initItemTable(), d)
	return
}

func (d *{{ .NodeName }}) initItemTable() []*adapter.Column {
	columns := make([]*adapter.Column, 0)
	var column *adapter.Column

	{{ range .PrimaryColumns }}
	column = &adapter.Column{
		Type:  "",
		Name:  "{{ .Name }}",
		Alias: "{{ .Alias }}",
	}
	column.SetPrimaryKey()
	columns = append(columns, column)
	{{ end }}

	{{ range .Columns }}
	columns = append(columns, &adapter.Column{
		Type:  "",
		Name:  "{{ .Name }}",
		Alias: "{{ .Alias }}",
	})
	{{ end }}

	return columns
}

func (d *{{ .NodeName }}) Name() string {
	return "{{ .PluralName }}"
}

func (d *{{ .NodeName }}) Type() core.FieldType {
	return FieldType{{ .NodeName }}
}

func (d *{{ .NodeName }}) IsList() bool {
	return true
}

func (d *{{ .NodeName }}) BuildFields() []*core.Field {
	fields := make([]*core.Field, 0)

	{{ range .Fields }}
	fields = append(fields,
		core.NewNodeField("{{ .Name }}", core.FieldType{{ .Type }}),
	)
	{{ end }}

	fields = append(fields, d.demoField())

	return fields
}

// DEMO: custom field and resolver
func (d *{{ .NodeName }}) demoField() *core.Field {
	field := core.NewNodeField("demo", core.FieldTypeString)
	field.SetResolver(func(p graphql.ResolveParams) (interface{}, error) {
		return "demo", nil
	})
	return field
}



`

type Parser struct {
	NodeName      string
	NodeNameLower string
	PluralName    string
	TableName     string

	PrimaryColumns []*Column
	Columns        []*Column
	Fields         []*Field
}

func NewParser(tableName string) *Parser {
	return &Parser{
		NodeName:      ToCamelCase(tableName),
		NodeNameLower: ToLowerCamelCase(tableName),
		PluralName:    Pluralize(tableName),
		TableName:     tableName,
	}
}

func (p *Parser) AddColumns(column *Column) {
	if column.IsPrimaryKey {
		p.PrimaryColumns = append(p.PrimaryColumns, column)
	} else {
		p.Columns = append(p.Columns, column)
	}

	p.Fields = append(p.Fields, &Field{
		Name: column.Name,
		Type: column.SwitchType(),
	})
	return
}

func (p *Parser) ParseTemplate() ([]byte, error) {
	// 解析模板
	tmpl, err := template.New("template").Parse(templateString) // templateString 是上面定义的模板字符串
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	// 渲染模板
	err = tmpl.Execute(buf, p)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Column struct {
	Name         string
	Alias        string
	Type         string
	IsPrimaryKey bool
}

func (c *Column) SwitchType() core.FieldType {
	switch strings.ToUpper(c.Type) {
	case "INT", "BIGINT", "SMALLINT", "TINYINT":
		return core.FieldTypeInt
	case "VARCHAR", "TEXT", "CHAR":
		return core.FieldTypeString
	case "DECIMAL":
		return core.FieldTypeFloat
	// ... 其他数据类型
	default:
		fmt.Printf("unknown data type: %s, please summit issue to https://github.com/shuishiyuanzhong/go-gql-builder/issues\n", c.Type)
		return "interface{}"
	}
}

type Field struct {
	Name string
	Type core.FieldType
}
