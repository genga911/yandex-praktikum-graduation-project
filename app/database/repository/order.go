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
	Accrual float32 `json:"accrual"`
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
	var accrual int
	err := or.DB.Connection.QueryRow(
		context.Background(),
		query,
		order.Number,
	).Scan(&order.UserID, &order.Status, &accrual, &order.UploadedAt)

	order.Accrual = float32(accrual) / 100

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	return &order, nil
}

func (or *Order) List(u *models.User, all bool) ([]*models.Order, error) {
	query := fmt.Sprintf("SELECT user_id, number, status, accrual, uploaded_at FROM %s WHERE user_id = $1", models.OrdersTableName)
	var args []interface{}
	args = append(args, u.ID)
	if !all {
		query += " AND status != $2"
		args = append(args, OrderStatusProcessed)
	}

	rows, err := or.DB.Connection.Query(
		context.Background(),
		query,
		args...,
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
		var accrual int
		err = rows.Scan(&order.UserID, &order.Number, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		order.Accrual = float32(accrual) / 100

		slice = append(slice, &order)
	}

	return slice, nil
}

// Update
func (or *Order) Update(o *models.Order) error {
	_, err := or.DB.Connection.Exec(
		context.Background(),
		fmt.Sprintf("UPDATE %s SET status=$1, accrual=$2", models.OrdersTableName),
		o.Status,
		int(o.Accrual*100),
	)

	return err
}

func (or *Order) GetFromAccrual(cfg *config.Config, o *models.Order) (*AccrualOrder, error) {
	// запрос к апи серверу за информацией по начислению баллов
	resp, err := http.Get(cfg.GetAccuralRequestAddress(o.Number))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не корректный код ответа от accural: %v", resp.StatusCode)
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

	// сохраним данные по заказу
	o.Status = aOrder.Status
	o.Accrual = aOrder.Accrual
	err = or.Update(o)
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

func (or *Order) Sync(cfg *config.Config, u *models.User) error {
	orders, err := or.List(u, false)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for index, order := range orders {
		wg.Add(1)
		go func(cfg *config.Config, order *models.Order, index int) {
			defer wg.Done()
			aOrder, err := or.GetFromAccrual(cfg, order)
			if err != nil {
				return
			}
			order.Status = aOrder.Status
			order.Accrual = aOrder.Accrual

			err = or.Update(order)
			if err != nil {
				panic(err)
			}
		}(cfg, order, index)
	}

	wg.Wait()

	return err
}

//GetBalance получение баланса
func (or *Order) GetBalanceSum(u *models.User) (float32, error) {
	var balance int
	err := or.DB.Connection.QueryRow(
		context.Background(),
		fmt.Sprintf("SELECT COALESCE(SUM(accrual), 0) FROM %s WHERE user_id=$1 AND status=$2", models.OrdersTableName),
		u.ID,
		OrderStatusProcessed,
	).Scan(&balance)

	return float32(balance) / 100, err
}
