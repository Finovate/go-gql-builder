package main

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/shuishiyuanzhong/go-gql-builder/example/user-department/conf"
	"github.com/shuishiyuanzhong/go-gql-builder/example/user-department/model"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"
)

func InitGraphQL() (*graphql.Schema, error) {
	core.DefaultRegistry().Register(model.NewUserDelegate())
	core.DefaultRegistry().Register(model.NewDepartmentDelegate())

	core.DefaultRegistry().SetDB(conf.C().Mysql.GetDB())
	return core.DefaultRegistry().BuildSchema()
}

func main() {
	// 定义Schema
	//schema := createSchema()

	schema, err := InitGraphQL()
	if err != nil {
		panic(err)
	}

	// 设置HTTP服务器
	httpHandler := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
	})
	http.Handle("/graphql", httpHandler)
	http.ListenAndServe(":8080", nil)
}
