package repository

import (
	"ddl2go/internal/entity"
)

type IDatabase interface {
	GetTables() []string
	GetColumns(tableName string) (columns []entity.Column)
	Close()
}
