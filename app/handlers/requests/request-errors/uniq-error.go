package request_errors

import "github.com/jackc/pgconn"

const OrderAlreadyExists = "Номер заказа уже был загружен этим пользователем"
const OrderCreatedByAnotherUser = "Номер заказа уже был загружен другим пользователем"

type UniqError struct {
	pgconn.PgError
}

func (ue *UniqError) Error() string {
	return ue.Message
}
