package entity

type (
	Table struct {
		Name    string
		Columns []Column
	}

	Column struct {
		Name        string
		GoType      string
		GraphQLType string
		GoPackage   string
		Comment     string
		IsComment   bool
		Tag         string
	}

	Reference struct {
		ColumnName    string
		RefTableName  string
		RefColumnName string
	}

	Row struct {
		ColumnName string `db:"column_name"`
		DataType   string `db:"data_type"`
		IsNullable string `db:"is_nullable"`
		Comment    string `db:"comment"`
		RefTable   string `db:"ref_table"`
		RefColumn  string `db:"ref_column"`
	}
)
