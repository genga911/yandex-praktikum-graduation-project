package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	"github.com/gin-gonic/gin"
	"github.com/theplant/luhn"
)

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Number string  `json:"order"`
	Sum    float32 `json:"sum"`
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

	orderNumber, _ := strconv.Atoi(request.Number)
	if !luhn.Valid(orderNumber) {
		c.AbortWithError(http.StatusUnprocessableEntity, errors.New("номер не корректен"))
		return nil
	}

	ro := repository.Order{
		DB: db,
	}
	rw := repository.Withdraw{
		DB: db,
	}
	balance, err := ro.GetBalanceSum(user)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	cw, err := rw.GetWithdrawSum(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	if (balance - cw) < request.Sum {
		c.AbortWithError(http.StatusPaymentRequired, errors.New("на счету не достаточно средств"))
		return nil
	}

	o := models.Order{Number: request.Number, UserID: user.ID}
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

func GetBalance(db *database.DB, c *gin.Context) *Balance {
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

	balance.Current, err = ro.GetBalanceSum(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}

	rw := repository.Withdraw{
		DB: db,
	}

	balance.Withdrawn, err = rw.GetWithdrawSum(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	balance.Current = balance.Current - balance.Withdrawn
	return &balance
}
