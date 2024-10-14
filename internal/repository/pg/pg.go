package repository_pg

import (
	"ddl2go/internal/entity"
	"ddl2go/internal/types"
	"ddl2go/internal/utils"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Source struct {
	db *sqlx.DB
}

func New(driver, connStr string) *Source {
	db, err := sqlx.Open(driver, connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &Source{
		db: db,
	}
}

func (s *Source) Close() {
	s.db.Close()
}

func (s *Source) GetTables() []string {
	rows, err := s.db.Query(`SELECT distinct table_name FROM information_schema.columns WHERE table_schema not in ('information_schema','pg_catalog')`)
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
	if rows.Err() != nil {
		log.Fatal(rows.Err())
	}

	return tables
}
func (s *Source) GetColumns(tableName string) (columns []entity.Column) {
	var rows []entity.Row
	err := s.db.Select(&rows, `
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
		graphqlType := types.SqlTypeToGraphQL(rows[i].DataType)
		goType, goPkg := types.SqlTypeToGo(rows[i].DataType, rows[i].IsNullable == "YES")
		fieldName, jsonTag := utils.ToCamelCase(rows[i].ColumnName)

		if rows[i].RefColumn != "" {
			rtn, _ := utils.ToCamelCase(rows[i].RefTable)
			rfn, _ := utils.ToCamelCase(rows[i].RefColumn)

			rows[i].Comment += rows[i].Comment + " (ref to " + rtn + "." + rfn + ")"
			graphqlType = "[" + rtn + "!]!"
		}

		b := &strings.Builder{}
		b.WriteString("`json:")
		b.WriteRune('"')
		b.WriteString(jsonTag)
		b.WriteRune('"')
		b.WriteString(" db:")
		b.WriteRune('"')
		b.WriteString(rows[i].ColumnName)
		b.WriteRune('"')
		b.WriteRune('`')

		columns = append(columns,
			entity.Column{
				Name:        fieldName,
				GoType:      goType,
				GraphQLType: graphqlType,
				GoPackage:   goPkg,
				Comment:     rows[i].Comment,
				IsComment:   rows[i].Comment != "",
				Tag:         b.String(),
			},
		)
	}
	return
}
