package requests

import (
	"errors"
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	"github.com/gin-gonic/gin"
)

type (
	LoginRequest struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

// DoRegister регистрация пользователя
func DoLogin(db *database.DB, c *gin.Context) *models.User {
	var request LoginRequest
	var err error
	err = c.ShouldBind(&request)
	if err != nil {
		// грубая ошибка валидации
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	rp := repository.User{
		DB: db,
	}

	var user *models.User
	// создадим пользователя, и проверим тип ошибки если такая возникла
	user, err = rp.GetUserByLoginPassword(request.Login, request.Password)

	// Если была ошибка, проверим её тип
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	if user == nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Неверная пара логин/пароль"))
	}

	return user
}
