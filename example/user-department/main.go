package main

import (
	"github.com/shuishiyuanzhong/go-gql-builder/example/user-department/conf"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/shuishiyuanzhong/go-gql-builder/example/user-department/model"
	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"
)

func InitGraphQL() (*graphql.Schema, error) {
	core.Registry().Register(model.NewUserDelegate())
	core.Registry().Register(model.NewDepartmentDelegate())

	core.Registry().SetDB(conf.C().Mysql.GetDB())
	return core.Registry().BuildSchema()
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
