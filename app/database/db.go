package database

import (
	"context"
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	Connection *pgx.Conn
}

// соединение с БД
func (db *DB) connect(connStr string) {
	connection, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(fmt.Sprintf("DB connection error: %s", err))
	}
	db.Connection = connection

	err = db.Connection.Ping(context.Background())
	if err != nil {
		panic(fmt.Sprintf("DB ping error: %s", err))
	}
}

// GetDB получить инстанс DB
func GetDB(cfg *config.Config) *DB {
	db := DB{}
	// подключимся к БД
	db.connect(cfg.DataBaseURI)

	// создать таблицы
	db.createTables()

	return &db
}

// создать таблицу пользователей
func (db *DB) createTables() {
	var tables []models.Model
	tables = append(tables, &models.User{})
	tables = append(tables, &models.Order{})
	tables = append(tables, &models.Withdraw{})

	for _, table := range tables {
		// дропаем старое
		_, err := db.Connection.Exec(context.Background(), table.DropTable())
		if err != nil {
			panic(fmt.Sprintf("Cannot drop table %s: %s", table.GetTableName(), err))
		}
		// создаем новое
		_, err = db.Connection.Exec(context.Background(), table.GetCreateTable())
		if err != nil {
			panic(fmt.Sprintf("Cannot create table %s: %s", table.GetTableName(), err))
		}
	}
}
