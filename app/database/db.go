package database

import (
	"context"
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	connection *pgx.Conn
}

// соединение с БД
func (db *DB) connect(connStr string) {
	connection, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(fmt.Sprintf("DB connection error: %s", err))
	}
	db.connection = connection

	err = db.connection.Ping(context.Background())
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

	for _, table := range tables {
		_, err := db.connection.Exec(context.Background(), table.GetCreateTable())
		if err != nil {
			panic(fmt.Sprintf("Cannot create table users: %s", err))
		}
	}
}
