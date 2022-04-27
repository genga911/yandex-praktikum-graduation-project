package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/handlers/requests/exceptions"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type AccrualOrder struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

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
		pgerr := exceptions.UniqError{}
		pgerr.Code = pgerrcode.UniqueViolation
		if order.UserID != u.ID {
			pgerr.Message = exceptions.OrderCreatedByAnotherUser
		} else {
			pgerr.Message = exceptions.OrderAlreadyExists
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

// Find получить заказ по номеру
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

func (or *Order) List(u *models.User) ([]*models.Order, error) {
	query := fmt.Sprintf("SELECT user_id, number, status, accrual, uploaded_at FROM %s WHERE user_id = $1", models.OrdersTableName)
	rows, err := or.DB.Connection.Query(
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

	var slice []*models.Order
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		slice = append(slice, &order)
	}

	return slice, nil
}

func (or *Order) GetFromAccrual(cfg *config.Config, o *models.Order) (*AccrualOrder, error) {
	// запрос к апи серверу за информацией по начислению баллов
	resp, err := http.Get(cfg.GetAccuralRequestAddress(o.Number))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// парсинг тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var aOrder AccrualOrder
	err = json.Unmarshal(body, &aOrder)
	if err != nil {
		return nil, err
	}

	return &aOrder, nil
}

func (or *Order) Delete(o *models.Order) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE number = $1", models.OrdersTableName)

	err := or.DB.Connection.QueryRow(
		context.Background(),
		query,
		o.Number,
	).Scan()

	return err
}

//GetBalance получение баланса
func (or *Order) GetBalance(cfg *config.Config, u *models.User) (float64, error) {
	orders, err := or.List(u)
	if err != nil {
		return 0, err
	}
	var wg sync.WaitGroup
	accruals := make([]*AccrualOrder, len(orders))

	for index, order := range orders {
		wg.Add(1)
		go func(cfg *config.Config, order *models.Order, index int) {
			defer wg.Done()
			aOrder, err := or.GetFromAccrual(cfg, order)
			if err != nil {
				panic(err)
			}

			accruals[index] = aOrder
		}(cfg, order, index)
	}

	wg.Wait()
	var balance float64
	for _, a := range accruals {
		if a.Status == OrderStatusProcessed {
			balance += a.Accrual
		}
	}

	return balance, err
}
