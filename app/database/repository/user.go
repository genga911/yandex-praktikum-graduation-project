package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/helpers"
	"github.com/jackc/pgx/v4"
)

type User struct {
	DB *database.DB
}

// создать пользователя
func (ur *User) Create(u *models.User) error {
	err := ur.DB.Connection.QueryRow(
		context.Background(),
		fmt.Sprintf("INSERT INTO %s(id, login, password, balance, withdraw) VALUES(DEFAULT, $1, $2, $3, $4) RETURNING id", models.UsersTableName),
		u.Login,
		helpers.MakeMD5(u.Password),
		u.Balance,
		u.Withdraw,
	).Scan(&u.ID)

	return err
}

// получить пользователя по иду
func (ur *User) Find(ID int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, login, balance, withdraw FROM %s WHERE id = $1 LIMIT 1", models.UsersTableName)

	err := ur.DB.Connection.QueryRow(
		context.Background(),
		query,
		ID,
	).Scan(&user.ID, &user.Login, &user.Balance, &user.Withdraw)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur *User) GetUserByLoginPassword(l string, p string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id, login, balance, withdraw FROM %s WHERE login = $1 AND password = $2 LIMIT 1", models.UsersTableName)

	err := ur.DB.Connection.QueryRow(
		context.Background(),
		query,
		l,
		helpers.MakeMD5(p),
	).Scan(&user.ID, &user.Login, &user.Balance, &user.Withdraw)

	// если есть ошибка и это не отсутствие результата
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur *User) IncreaseBalance(sum float64, u *models.User) error {
	query := fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE id = $2 RETURNING balance", models.UsersTableName)
	err := ur.DB.Connection.QueryRow(
		context.Background(),
		query,
		sum,
		u.ID,
	).Scan(&u.Balance)

	return err
}

func (ur *User) IncreaseWithdraw(sum float64, u *models.User) error {
	query := fmt.Sprintf("UPDATE %s SET withdraw = withdraw + $1 WHERE id = $2 RETURNING withdraw", models.UsersTableName)
	err := ur.DB.Connection.QueryRow(
		context.Background(),
		query,
		sum,
		u.ID,
	).Scan(&u.Withdraw)

	return err
}
