package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	request_errors "github.com/genga911/yandex-praktikum-graduation-project/app/handlers/requests/exceptions"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type Order struct {
	DB *database.DB
}

const OrderStatusNew = "NEW"
const OrderStatusProcessing = "PROCESSING"
const OrderStatusInvalid = "INVALID"
const OrderStatusProcessed = "PROCESSED"

// Create создать заказ
func (or *Order) Create(u *models.User, o *models.Order) error {
	// проверим, не создавали ли пользователи ранее этот заказ
	order, err := or.Find(o.Number)

	if err != nil {
		return err
	}

	if order != nil {
		// если заказ создан
		// сымитируем ошибку уникальности pgconn
		pgerr := request_errors.UniqError{}
		pgerr.Code = pgerrcode.UniqueViolation
		if order.UserID != u.ID {
			pgerr.Message = request_errors.OrderCreatedByAnotherUser
		} else {
			pgerr.Message = request_errors.OrderAlreadyExists
		}

		o = order
		return &pgerr
	}

	err = or.DB.Connection.QueryRow(
		context.Background(),
		fmt.Sprintf("INSERT INTO %s(number, user_id) VALUES($1, $2) RETURNING uploaded_at", models.OrdersTableName),
		o.Number,
		u.ID,
	).Scan(&o.UploadedAt)

	// пропишем дефолтные значения
	o.Accrual = 0
	o.Status = OrderStatusNew

	return err
}

// получить заказ по номеру
func (or *Order) Find(Number string) (*models.Order, error) {
	order := models.Order{Number: Number}
	query := fmt.Sprintf("SELECT user_id, status, accrual, uploaded_at FROM %s WHERE number = $1 LIMIT 1", models.OrdersTableName)

	err := or.DB.Connection.QueryRow(
		context.Background(),
		query,
		order.Number,
	).Scan(&order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	return &order, nil
}
