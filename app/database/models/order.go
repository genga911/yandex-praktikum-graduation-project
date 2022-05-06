package models

import (
	"fmt"
	"time"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
	UserID     int       `json:"-"`
}

const OrdersTableName = "orders"

func (o *Order) GetCreateTable() string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"number varchar(255) not null UNIQUE,"+
			"user_id int not null,"+
			"status varchar(255) not null default 'NEW',"+
			"accrual int not null default 0,"+
			"uploaded_at timestamp not null default CURRENT_TIMESTAMP"+
			");", OrdersTableName)
}

func (o *Order) DropTable() string {
	return fmt.Sprintf("DROP table IF EXISTS %s; DROP TYPE IF EXISTS order_status;", OrdersTableName)
}

func (o *Order) GetTableName() string {
	return OrdersTableName
}
