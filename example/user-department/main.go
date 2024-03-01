package main

import (
	"net/http"

	"github.com/shuishiyuanzhong/go-gql-builder/pkg/core"

	"github.com/shuishiyuanzhong/go-gql-builder/example/user-department/conf"
	"github.com/shuishiyuanzhong/go-gql-builder/example/user-department/model"
)

func InitGraphQL() (http.Handler, error) {
	core.DefaultRegistry().Register(model.NewUserDelegate())
	core.DefaultRegistry().Register(model.NewDepartmentDelegate())

	core.DefaultRegistry().SetDB(conf.C().Mysql.GetDB())
	return core.DefaultRegistry().BuildHandler()
}

func main() {
	// 定义Schema
	//schema := createSchema()

	graphqlHandler, err := InitGraphQL()
	if err != nil {
		panic(err)
	}

	http.Handle("/graphql", graphqlHandler)
	http.ListenAndServe(":8080", nil)
}
