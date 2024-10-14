package main

import (
	"ddl2go/internal/entity"
	"ddl2go/internal/generators"
	"ddl2go/internal/repository"
	repository_ms "ddl2go/internal/repository/ms"
	repository_pg "ddl2go/internal/repository/pg"
	"ddl2go/internal/utils"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

type (
	Env struct {
		Package    string
		ConnStr    string
		Driver     string
		GenGraphQL bool
	}
)

func main() {
	env := Env{
		Package:    os.Getenv("PACKAGE"),
		ConnStr:    os.Getenv("CONN_STR"),
		Driver:     os.Getenv("DRIVER"),
		GenGraphQL: os.Getenv("GRAPHQL") == "true",
	}
	_ = os.MkdirAll(env.Package, os.ModePerm)

	var db repository.IDatabase
	switch env.Driver {
	case "postgres":
		db = repository_pg.New(env.Driver, env.ConnStr)
	case "sqlserver":
		db = repository_ms.New(env.Driver, env.ConnStr)
	}
	defer db.Close()

	tables := db.GetTables()
	for i := range tables {
		columns := db.GetColumns(tables[i])

		goStructName, _ := utils.ToCamelCase(tables[i])
		table := entity.Table{
			Name:    goStructName,
			Columns: columns,
		}
		generators.GoStruct(env.Package, tables[i], table)
		if env.GenGraphQL {
			generators.GraphQLSchema(env.Package, tables[i], table)
		}
	}
	fmt.Println("completed")
}
