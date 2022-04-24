package models

import "fmt"

type User struct {
	ID       int     `json:"id"`
	Login    string  `json:"login"`
	Password string  `json:"-"`
	Balance  float64 `json:"balance,omitempty"`
	Withdraw float64 `json:"withdraw"`
}

const UsersTableName = "users"

func (ur *User) GetCreateTable() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"id serial not null,"+
		"login varchar(255) not null UNIQUE,"+
		"password varchar(255) not null,"+
		"balance NUMERIC(3,2) not null default 0,"+
		"withdraw NUMERIC(3,2) not null default 0"+
		");", UsersTableName)
}

func (ur *User) DropTable() string {
	return fmt.Sprintf("DROP table IF EXISTS %s;", UsersTableName)
}

func (ur *User) GetTableName() string {
	return UsersTableName
}
