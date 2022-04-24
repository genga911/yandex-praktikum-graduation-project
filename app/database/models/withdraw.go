package models

import (
	"fmt"
	"time"
)

type Withdraw struct {
	Number      string    `json:"number"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      int       `json:"-"`
}

const WithdrawnTableName = "withdrawals"

func (o *Withdraw) GetCreateTable() string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ("+
			"number varchar(255) not null UNIQUE,"+
			"sum int not null,"+
			"processed_at timestamp not null default CURRENT_TIMESTAMP,"+
			"user_id int not null"+
			");", WithdrawnTableName)
}

func (o *Withdraw) DropTable() string {
	return fmt.Sprintf("DROP table IF EXISTS %s;", WithdrawnTableName)
}

func (o *Withdraw) GetTableName() string {
	return WithdrawnTableName
}
