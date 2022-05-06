package models

type Model interface {
	GetCreateTable() string
	DropTable() string
	GetTableName() string
}
