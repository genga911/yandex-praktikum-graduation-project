package requests

import (
	"errors"
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type (
	RegisterRequest struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

// DoRegister регистрация пользователя
func DoRegister(db *database.DB, c *gin.Context) *models.User {
	var request RegisterRequest
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

	// создадим пользователя, и проверим тип ошибки если такая возникла
	newUser := models.User{
		Login:    request.Login,
		Password: request.Password,
	}
	err = rp.Create(&newUser)

	// Если была ошибка, проверим её тип
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				c.AbortWithError(http.StatusConflict, errors.New("логин уже занят"))
				return nil
			}
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
			return nil
		}
	}

	return &newUser
}
