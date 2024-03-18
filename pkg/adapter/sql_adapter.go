package adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"

	sqlArgument "github.com/shuishiyuanzhong/go-gql-builder/pkg/adapter/internal/argument"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"
	coreArgument "github.com/shuishiyuanzhong/go-gql-builder/pkg/core/argument"
)

// SqlAdapter is a part of Node interface, which is
// designed to bridge business objects with SQL queries.
type SqlAdapter interface {
	Resolve() graphql.FieldResolveFn
}

// DefaultSqlAdapter is a default implementation of SqlAdapter.
// It is use to query single table with custom fields.
type DefaultSqlAdapter struct {
	node core.Node

	tableName      string
	tableColumns   []*Column
	columnsByAlias map[string]*Column
	columnsByName  map[string]*Column
	primaryKeys    []*Column
}

func NewDefaultSqlAdapter(tableName string, columns []*Column, node core.Node) *DefaultSqlAdapter {
	d := &DefaultSqlAdapter{
		node:           node,
		tableName:      tableName,
		tableColumns:   make([]*Column, 0, len(columns)),
		columnsByAlias: make(map[string]*Column),
		columnsByName:  make(map[string]*Column),
		primaryKeys:    make([]*Column, 0),
	}

	for _, column := range columns {
		d.tableColumns = append(d.tableColumns, column)
		d.columnsByAlias[column.Alias] = column
		d.columnsByName[column.Name] = column
		if column.IsPrimaryKey() {
			d.primaryKeys = append(d.primaryKeys, column)
		}
	}

	return d
}

func (d *DefaultSqlAdapter) Resolve() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		customFields := make([]*ast.Field, 0)
		for _, field := range p.Info.FieldASTs {
			if field.Name.Value == d.node.Name() {
				selections := field.SelectionSet.Selections
				for _, selection := range selections {
					customFields = append(customFields, selection.(*ast.Field))
				}
				break
			}
		}

		var customCollect []string
		for _, field := range customFields {
			if _, ok := d.columnsByName[field.Name.Value]; ok {
				customCollect = append(customCollect, field.Name.Value)
			}
		}

		if len(customCollect) <= 0 {
			for _, pk := range d.primaryKeys {
				customCollect = append(customCollect, pk.Name)
			}
			if len(customCollect) == 0 {
				customCollect = []string{"*"}
			}
		}

		qc := sqlArgument.NewQueryClauses(strings.Join(customCollect, ","), d.tableName)

		for name, value := range p.Args {
			arg := coreArgument.Factory(name)
			if arg == nil {
				return nil, fmt.Errorf("argument typename is not exist")
			}
			err := arg.Validate(value)
			if err != nil {
				return nil, err
			}

			sqlArg, ok := arg.(sqlArgument.SqlArgument)
			if ok {
				sqlArg.CombineSql(qc)
				//switch x := sqlArg.(type) {
				//case *argument.FilterArgument:
				//	qc.SetWhere(x.ParseSqlValue())
				//}
			}
		}

		sql, err := qc.ToSql()
		if err != nil {
			return nil, err
		}

		rows, err := d.node.GetRegistry().GetDB().QueryContext(context.Background(), sql)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		columns, _ := rows.Columns()

		cache := make([]interface{}, len(columns)) // 临时存储每行数据
		for i := range cache {                     // 为每一列初始化一个指针
			var a interface{}
			cache[i] = &a
		}
		var list []map[string]interface{} //返回的切片
		for rows.Next() {
			_ = rows.Scan(cache...)

			item := make(map[string]interface{})
			for i, col := range columns {
				val := *(cache[i].(*interface{})) // 获取实际类型的值

				if bytesVal, ok := val.([]byte); ok {
					val = string(bytesVal)
				}
				item[col] = val
			}
			fmt.Println(item)
			list = append(list, item)
		}

		return list, nil
		//return []map[string]interface{}{{"id": "1", "name": "Example Product", "price": 99.99}}, nil
	}
}

// DTO to GraphQL Object

type Column struct {
	Type         ColumnType
	Name         string
	Alias        string
	isPrimaryKey bool
}

func (c *Column) SetPrimaryKey() {
	c.isPrimaryKey = true
}

func (c *Column) IsPrimaryKey() bool {
	return c.isPrimaryKey
}

type ColumnType string

const (
	Int    ColumnType = "Int"
	Float  ColumnType = "Float"
	String ColumnType = "String"
)
