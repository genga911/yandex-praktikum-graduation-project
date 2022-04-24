package requests

import (
	"errors"
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	"github.com/gin-gonic/gin"
)

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type WithdrawRequest struct {
	Number string  `json:"order" binding:"required"`
	Sum    float32 `json:"sum" binding:"required"`
}

// RegisterWithdraw регистрация списания
func RegisterWithdraw(db *database.DB, c *gin.Context) *models.Withdraw {
	u, exist := c.Get("user")
	if !exist {
		c.AbortWithError(http.StatusUnauthorized, errors.New("пользователь не найден"))
		return nil
	}

	user := u.(*models.User)

	var request WithdrawRequest
	var err error
	err = c.ShouldBind(&request)
	if err != nil {
		// грубая ошибка валидации
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	if user.Balance < request.Sum {
		c.AbortWithError(http.StatusPaymentRequired, errors.New("на счету не достаточно средств"))
		return nil
	}

	rw := repository.Withdraw{
		DB: db,
	}

	o := models.Order{Number: request.Number}
	wi, err := rw.Create(request.Sum, &o)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	return wi
}

func ListWithdraw(db *database.DB, c *gin.Context) []*models.Withdraw {
	u, exist := c.Get("user")
	if !exist {
		c.AbortWithError(http.StatusUnauthorized, errors.New("пользователь не найден"))
		return nil
	}
	user := u.(*models.User)

	r := repository.Withdraw{
		DB: db,
	}

	list, err := r.List(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	return list
}
