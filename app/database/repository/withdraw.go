package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/jackc/pgx/v4"
)

type Withdraw struct {
	DB *database.DB
}

// Create создать заказ
func (w *Withdraw) Create(sum float32, o *models.Order) (*models.Withdraw, error) {
	var wi models.Withdraw
	err := w.DB.Connection.QueryRow(
		context.Background(),
		fmt.Sprintf("INSERT INTO %s(number, sum, user_id) VALUES($1, $2, $3) RETURNING processed_at", models.WithdrawnTableName),
		o.Number,
		int(sum*100),
		o.UserID,
	).Scan(&wi.ProcessedAt)

	wi.Number = o.Number
	wi.Sum = sum

	if err != nil {
		return nil, err
	}

	return &wi, nil
}

// List лист пользователя со списаниями
func (w *Withdraw) List(u *models.User) ([]*models.Withdraw, error) {
	query := fmt.Sprintf("SELECT number, sum, processed_at, user_id FROM %s WHERE user_id = $1 ORDER BY processed_at DESC", models.WithdrawnTableName)
	rows, err := w.DB.Connection.Query(
		context.Background(),
		query,
		u.ID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	defer rows.Close()

	var slice []*models.Withdraw
	for rows.Next() {
		var withdraw models.Withdraw
		err = rows.Scan(&withdraw.Number, &withdraw.Sum, &withdraw.ProcessedAt, &withdraw.UserID)
		if err != nil {
			return nil, err
		}

		slice = append(slice, &withdraw)
	}

	return slice, nil
}

func (w *Withdraw) GetWithdrawSum(u *models.User) (float32, error) {
	var withdraw int

	err := w.DB.Connection.QueryRow(
		context.Background(),
		fmt.Sprintf("SELECT COALESCE(SUM(sum), 0) FROM %s WHERE user_id=$1", models.WithdrawnTableName),
		u.ID,
	).Scan(&withdraw)

	return float32(withdraw) / 100, err
}
