package repository_ms

import (
	"ddl2go/internal/entity"
	"ddl2go/internal/types"
	"ddl2go/internal/utils"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	rows, err := s.db.Query(`SELECT DISTINCT table_name FROM information_schema.columns WHERE table_schema NOT IN ('information_schema', 'sys') AND table_name not in ('spt_values')`)
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
        cp.name AS "column",
        tr.name AS "ref_table",
        cr.name AS "ref_column"
    FROM sys.foreign_keys AS fk
    INNER JOIN sys.foreign_key_columns AS fkc ON fk.object_id = fkc.constraint_object_id
    INNER JOIN sys.tables AS tp ON fkc.parent_object_id = tp.object_id
    INNER JOIN sys.columns AS cp ON fkc.parent_object_id = cp.object_id AND fkc.parent_column_id = cp.column_id
    INNER JOIN sys.tables AS tr ON fkc.referenced_object_id = tr.object_id
    INNER JOIN sys.columns AS cr ON fkc.referenced_object_id = cr.object_id AND fkc.referenced_column_id = cr.column_id
    WHERE 
        tp.name = @p1
)
SELECT
	cols.column_name, 
	cols.data_type, 
	cols.is_nullable,
	coalesce(
        (
            SELECT
                ep.value
            FROM sys.tables t
            INNER JOIN sys.columns c ON t.object_id = c.object_id
            LEFT JOIN sys.extended_properties ep ON ep.major_id = c.object_id AND ep.minor_id = c.column_id AND ep.name = 'MS_Description'
            WHERE 
                t.name = cols.table_name
                AND c.name = cols.column_name
        ),
    '') AS "comment",
    COALESCE(r.ref_table, '') as ref_table,
    COALESCE(r.ref_column, '') as ref_column
FROM information_schema.columns cols
LEFT JOIN REFS r on r."column" = cols.column_name
WHERE 
	cols.table_name = @p1`, tableName)
	if err != nil {
		log.Fatal(errors.Wrap(err, tableName))
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
