package adapter

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"
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
			// FIXME 这个判断的逻辑不太对，当进到Resolve方法的时候，p.Info.FieldASTs 按道理来说只有一个元素
			// 如： users{id,name,age}，实际上拿到的就是一个users
			// 下面的判断也不正确，tableName并不等于这个Node的name
			if field.Name.Value == d.tableName {
				selections := field.SelectionSet.Selections
				for _, selection := range selections {
					customFields = append(customFields, selection.(*ast.Field))
				}
				break
			}
		}

		// FIXME 判断有点问题，理应是<=
		if len(customFields) < 0 {
			return nil, errors.New("no custom fields to query")
		}

		var customCollect []string
		for _, field := range customFields {
			customCollect = append(customCollect, field.Name.Value)
		}
		// FIXME 这里需要额外加一些判断，比如graphql传递了users{id,name,department}
		// 但实际上user表只有id，name两个字段，department这个field并不存在于user表中，而是有额外的逻辑进行处理
		// 因此customCollect应该和tableColumns进行比较，如果里面的field不在tableColumns中，则剔除出去不进行查询

		sql := "SELECT %s from %s"
		sql = fmt.Sprintf(sql, strings.Join(customCollect, ","), d.tableName)

		rows, err := core.Registry().GetDB().QueryContext(context.Background(), sql)
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
