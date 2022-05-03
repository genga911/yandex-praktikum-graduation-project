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
func (db *DB) connect(connStr string) error {
	connection, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("DB connection error: %s", err)
	}
	db.Connection = connection

	err = db.Connection.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("DB ping error: %s", err)
	}

	return nil
}

// GetDB получить инстанс DB
func GetDB(cfg *config.Config) (*DB, error) {
	db := DB{}
	// подключимся к БД
	err := db.connect(cfg.DataBaseURI)
	if err != nil {
		return nil, err
	}

	// создать таблицы
	err = db.createTables()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// создать таблицу пользователей
func (db *DB) createTables() error {
	var tables []models.Model
	tables = append(tables, &models.User{})
	tables = append(tables, &models.Order{})
	tables = append(tables, &models.Withdraw{})

	for _, table := range tables {
		// создаем новое
		_, err := db.Connection.Exec(context.Background(), table.GetCreateTable())
		if err != nil {
			return fmt.Errorf("cannot create table %s: %s", table.GetTableName(), err)
		}
	}

	return nil
}
