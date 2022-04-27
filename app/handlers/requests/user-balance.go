package requests

import (
	"errors"
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	"github.com/gin-gonic/gin"
)

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Number string  `json:"order"`
	Sum    float64 `json:"sum"`
}

// RegisterWithdraw регистрация списания
func RegisterWithdraw(cfg *config.Config, db *database.DB, c *gin.Context) *models.Withdraw {
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

	ro := repository.Order{
		DB: db,
	}

	balance, err := ro.GetBalance(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	if balance < request.Sum {
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

	rw := repository.Withdraw{
		DB: db,
	}

	list, err := rw.List(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	return list
}

func GetBalance(cfg *config.Config, db *database.DB, c *gin.Context) *Balance {
	u, exist := c.Get("user")
	if !exist {
		c.AbortWithError(http.StatusUnauthorized, errors.New("пользователь не найден"))
		return nil
	}
	user := u.(*models.User)
	var balance Balance
	ro := repository.Order{
		DB: db,
	}
	var err error
	balance.Current, err = ro.GetBalance(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	rw := repository.Withdraw{
		DB: db,
	}

	list, err := rw.List(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	balance.Withdrawn = 0
	for _, w := range list {
		balance.Withdrawn += w.Sum
	}

	return &balance
}
