package exceptions

import "github.com/jackc/pgconn"

const OrderAlreadyExists = "номер заказа уже был загружен этим пользователем"
const OrderCreatedByAnotherUser = "номер заказа уже был загружен другим пользователем"

type UniqError struct {
	pgconn.PgError
}

func (ue *UniqError) Error() string {
	return ue.Message
}
