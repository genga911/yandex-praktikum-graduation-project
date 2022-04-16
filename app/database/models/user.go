package models

type User struct {
	ID      string
	Name    string
	Balance int
}

func (u *User) GetCreateTable() string {
	return "CREATE TABLE IF NOT EXISTS users (" +
		"id serial not null," +
		"login varchar(255) not null," +
		"password varchar(255) not null," +
		"balance int not null default 0" +
		");"
}
