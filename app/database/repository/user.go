package repository

import (
	"context"
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/helpers"
)

type User struct {
	DB *database.DB
}

// создать пользователя
func (ur User) Create(u *models.User) error {
	err := ur.DB.Connection.QueryRow(
		context.Background(),
		fmt.Sprintf("INSERT INTO %s(id, login, password, balance) VALUES(DEFAULT, $1, $2, $3) RETURNING id", models.UsersTableName),
		u.Login,
		helpers.MakeMD5(u.Password),
		u.Balance,
	).Scan(&u.ID)

	return err
}

// получить пользователя по иду
func (ur User) Find(id int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT login, balance FROM %s WHERE id = $1 LIMIT 1", models.UsersTableName)

	err := ur.DB.Connection.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(user.Login, user.Balance)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
