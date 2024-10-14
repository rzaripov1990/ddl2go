package types

import "strings"

func SqlTypeToGo(sqlType string, isNull bool) (_type, _pkg string) {
	_allowNull := true
	switch strings.ToLower(sqlType) {
	case "serial", "integer", "smallint", "int":
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
	case "char", "varchar", "text", "character", "character varying", "nvarchar":
		_type = "string"
		_pkg = ""
	case "boolean":
		_type = "bool"
		_pkg = ""
	case "uuid":
		_type = "uuid.UUID"
		_pkg = "github.com/google/uuid"
	case "timestamp", "timestamptz", "date", "time", "datetime", "time without time zone", "timestamp with time zone", "timestamp without time zone":
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

func SqlTypeToGraphQL(sqlType string) string {
	switch strings.ToLower(sqlType) {
	case "serial", "integer", "smallint", "bigint":
		return "Int"
	case "decimal", "numeric", "money", "real", "double precision":
		return "Float"
	case "char", "varchar", "text", "character", "character varying":
		return "String"
	case "boolean":
		return "Boolean"
	case "uuid":
		return "ID"
	case "timestamp", "timestamptz", "date", "time", "time without time zone", "timestamp with time zone", "timestamp without time zone":
		return "String"
	case "bytea", "json", "jsonb", "xml":
		return "String"
	default:
		return "String"
	}
}
