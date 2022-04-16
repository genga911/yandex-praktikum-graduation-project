package repository

import (
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
)

type User struct {
	DB database.DB
}

// создать пользователя
func (ur User) Create(u *models.User) {

}

// создать пользователя
func (ur User) Find(id string) *models.User {
	return &models.User{}
}
