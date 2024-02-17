package adapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"
	"strings"
)

// SqlAdapter is a part of Node interface, which is
// designed to bridge business objects with SQL queries.
type SqlAdapter interface {
	Resolve() graphql.FieldResolveFn
}

// DefaultSqlAdapter is a default implementation of SqlAdapter.
// It is use to query single table with custom fields.
type DefaultSqlAdapter struct {
	tableName      string
	tableColumns   []*Column
	columnsByAlias map[string]*Column
	columnsByName  map[string]*Column
}

func NewDefaultSqlAdapter(tableName string, columns []*Column) *DefaultSqlAdapter {
	d := &DefaultSqlAdapter{
		tableName:      tableName,
		tableColumns:   make([]*Column, 0, len(columns)),
		columnsByAlias: make(map[string]*Column),
		columnsByName:  make(map[string]*Column),
	}

	for _, column := range columns {
		d.tableColumns = append(d.tableColumns, column)
		d.columnsByAlias[column.Alias] = column
		d.columnsByName[column.Name] = column
	}

	return d
}

func (d *DefaultSqlAdapter) Resolve() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {

		customFields := make([]*ast.Field, 0)
		for _, field := range p.Info.FieldASTs {
			if field.Name.Value == d.tableName {
				selections := field.SelectionSet.Selections
				for _, selection := range selections {
					customFields = append(customFields, selection.(*ast.Field))
				}
				break
			}
		}

		if len(customFields) < 0 {
			return nil, errors.New("no custom fields to query")
		}

		var customCollect []string
		for _, field := range customFields {
			customCollect = append(customCollect, field.Name.Value)
		}

		sql := "SELECT %s from %s"
		sql = fmt.Sprintf(sql, strings.Join(customCollect, ","), d.tableName)

		rows, err := core.Registry().GetDB().QueryContext(context.Background(), sql)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		columns, _ := rows.Columns()

		cache := make([]interface{}, len(columns)) //临时存储每行数据
		for index, _ := range cache {              //为每一列初始化一个指针
			var a interface{}
			cache[index] = &a
		}
		var list []map[string]interface{} //返回的切片
		for rows.Next() {
			_ = rows.Scan(cache...)

			item := make(map[string]interface{})
			for i, data := range cache {
				item[columns[i]] = *data.(*interface{}) //取实际类型
			}
			list = append(list, item)
		}

		return list, nil
		//return []map[string]interface{}{{"id": "1", "name": "Example Product", "price": 99.99}}, nil
	}
}

// DTO to GraphQL Object

type Column struct {
	Type  ColumnType
	Name  string
	Alias string
}

type ColumnType string

const (
	Int    ColumnType = "Int"
	Float  ColumnType = "Float"
	String ColumnType = "String"
)
