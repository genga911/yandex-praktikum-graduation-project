package models

import "fmt"

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Balance  int    `json:"balance,omitempty"`
}

const UsersTableName = "users"

func (u *User) GetCreateTable() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"id serial not null,"+
		"login varchar(255) not null UNIQUE,"+
		"password varchar(255) not null,"+
		"balance int not null default 0"+
		");", UsersTableName)
}

func (ur *User) DropTable() string {
	return fmt.Sprintf("DROP table IF EXISTS %s;", UsersTableName)
}

func (ur *User) GetTableName() string {
	return UsersTableName
}
