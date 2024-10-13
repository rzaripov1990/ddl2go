package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	_ "github.com/joho/godotenv/autoload"
)

type (
	Env struct {
		Package    string
		PostgresCS string
	}

	Table struct {
		Name    string
		Columns []Column
	}

	Column struct {
		Name      string
		Type      string
		GoPackage string
		Comment   string
		IsComment bool
		tagJson   string
		tagDB     string
		Tag       string
	}

	Reference struct {
		ColumnName    string
		RefTableName  string
		RefColumnName string
	}
)

func main() {
	env := Env{
		Package:    os.Getenv("PACKAGE"),
		PostgresCS: os.Getenv("PG_CONN_STR"),
	}
	_ = os.MkdirAll(env.Package, os.ModePerm)

	db, err := sqlx.Open("postgres", env.PostgresCS)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables := getTables(db)
	for i := range tables {
		columns := getColumns(db, tables[i])

		goStructName, _ := toCamelCase(tables[i])
		table := Table{
			Name:    goStructName,
			Columns: columns,
		}

		tmpl := `type {{ .Name }} struct {
{{- range .Columns }} 
	{{ .Name }} {{ .Type }} {{ .Tag }} {{if .IsComment }}//{{ .Comment }} {{ end }}
{{- end }}
}

type {{ .Name }}Arr []{{ .Name }}`
		t, err := template.New("struct").Parse(tmpl)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(env.Package + "/" + tables[i] + ".go")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		uniqPackages := map[string]bool{}
		imports := ""
		for i := range columns {
			if _, ok := uniqPackages[columns[i].GoPackage]; !ok && columns[i].GoPackage != "" {
				uniqPackages[columns[i].GoPackage] = true
				if imports == "" {
					imports += "import ("
				}
				imports += "\n    \"" + columns[i].GoPackage + "\""
			}
		}

		_, _ = file.WriteString("package " + env.Package + "\n\n")
		if imports != "" {
			_, _ = file.WriteString(imports + "\n)\n\n")
		}

		err = t.Execute(file, table)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("completed")
}

func getTables(db *sqlx.DB) []string {
	rows, err := db.Query(`SELECT distinct table_name FROM information_schema.columns WHERE table_schema not in ('information_schema','pg_catalog')`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		table  string
		tables []string
	)

	for rows.Next() {
		err := rows.Scan(&table)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, table)
	}

	return tables
}

type Row struct {
	ColumnName string `db:"column_name"`
	DataType   string `db:"data_type"`
	IsNullable string `db:"is_nullable"`
	Comment    string `db:"comment"`
	RefTable   string `db:"ref_table"`
	RefColumn  string `db:"ref_column"`
}

func getColumns(db *sqlx.DB, tableName string) (columns []Column) {
	var rows []Row
	err := db.Select(&rows, `
WITH REFS AS (
	SELECT 
		kcu1.column_name as column,
		kcu2.table_name AS ref_table,
		kcu2.column_name AS ref_column
	FROM information_schema.referential_constraints AS rc
	JOIN information_schema.key_column_usage AS kcu1 ON kcu1.constraint_name = rc.constraint_name AND kcu1.constraint_schema = rc.constraint_schema
	JOIN information_schema.key_column_usage AS kcu2 ON kcu2.constraint_name = rc.unique_constraint_name AND kcu2.constraint_schema = rc.constraint_schema AND kcu2.ordinal_position = kcu1.position_in_unique_constraint
	WHERE 
		kcu1.table_name = $1
) 
SELECT 
	cols.column_name, 
	cols.data_type, 
	cols.is_nullable,
	COALESCE(
		(
			SELECT
				pg_catalog.col_description(c.oid, cols.ordinal_position::int)
			FROM
				pg_catalog.pg_class c
			WHERE
				c.oid = (SELECT ('"' || cols.table_name || '"')::regclass::oid)
				AND c.relname = cols.table_name
		), 
	'') AS comment,
    COALESCE(r.ref_table, '') as ref_table,
    COALESCE(r.ref_column, '') as ref_column
FROM information_schema.columns cols 
LEFT JOIN REFS r ON r.column = cols.column_name
WHERE 
	cols.table_name = $1`, tableName)
	if err != nil {
		log.Fatal(err)
	}

	for i := range rows {
		goType, goPkg := sqlTypeToGo(rows[i].DataType, rows[i].IsNullable == "YES")
		fieldName, jsonTag := toCamelCase(rows[i].ColumnName)

		if rows[i].RefColumn != "" {
			rtn, _ := toCamelCase(rows[i].RefTable)
			rfn, _ := toCamelCase(rows[i].RefColumn)

			rows[i].Comment += rows[i].Comment + " (ref to " + rtn + "." + rfn + ")"
		}

		columns = append(columns,
			Column{
				Name:      fieldName,
				Type:      goType,
				GoPackage: goPkg,
				Comment:   rows[i].Comment,
				IsComment: rows[i].Comment != "",
				tagJson:   jsonTag,
				tagDB:     rows[i].ColumnName,
				Tag:       "`json:\"" + jsonTag + "\" db:\"" + rows[i].ColumnName + "\"`",
			},
		)
	}
	return
}

func sqlTypeToGo(sqlType string, isNull bool) (_type, _pkg string) {
	_allowNull := true
	switch strings.ToLower(sqlType) {
	case "serial", "integer", "smallint":
		_type = "int"
		_pkg = ""
	case "bigint":
		_type = "int64"
		_pkg = ""
	case "decimal", "numeric", "money":
		_type = "float64"
		_pkg = ""
	case "real", "double precision":
		_type = "float32"
		_pkg = ""
	case "char", "varchar", "text", "character", "character varying":
		_type = "string"
		_pkg = ""
	case "boolean":
		_type = "bool"
		_pkg = ""
	case "uuid":
		_type = "uuid.UUID"
		_pkg = "github.com/google/uuid"
	case "timestamp", "timestamptz", "date", "time", "time without time zone", "timestamp with time zone", "timestamp without time zone":
		_type = "time.Time"
		_pkg = "time"
	case "bytea", "json", "jsonb", "xml": //"inet", "cidr", "point", "line", "lseg", "box"
		_type = "[]byte"
		_pkg = ""
		_allowNull = false
	case "interval":
		_type = "time.Duration"
		_pkg = "time"
	case "int[]", "integer[]":
		_type = "[]int64"
		_pkg = ""
		_allowNull = false
	case "text[]":
		_type = "[]string"
		_pkg = ""
		_allowNull = false
	default:
		_type = "any"
		_pkg = ""
		_allowNull = false
	}
	if isNull && _allowNull {
		_type = "*" + _type
	}
	return
}

func toCamelCase(input string) (upperCS, classicCS string) {
	words := strings.Split(input, "_")
	for i := range words {
		lower := strings.ToLower(words[i][1:])
		upper := strings.ToUpper(words[i][:1])

		words[i] = upper + lower
	}
	upperCS = strings.Join(words, "")
	classicCS = strings.ToLower(upperCS[:1]) + upperCS[1:]
	return
}
