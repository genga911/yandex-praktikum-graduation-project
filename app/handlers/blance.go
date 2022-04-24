package handlers

import (
	"errors"
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/handlers/requests"
	"github.com/gin-gonic/gin"
)

type BalanceResponse struct {
	Balance  float64 `json:"balance"`
	Withdraw float64 `json:"withdraw"`
}

// Balance Баланс пользователя
func (uh *UserHandlers) Balance(c *gin.Context) {
	u, exist := c.Get("user")
	if !exist {
		c.AbortWithError(http.StatusUnauthorized, errors.New("пользователь не найден"))
		return
	}

	user := u.(*models.User)
	c.JSON(http.StatusOK, BalanceResponse{Balance: user.Balance, Withdraw: user.Withdraw})
}

func (uh *UserHandlers) RegisterWithdraw(c *gin.Context) {
	withdraw := requests.RegisterWithdraw(uh.DB, c)

	if withdraw != nil {
		c.JSON(http.StatusOK, withdraw)
	}
}

// ListWithdraw лист списаний
func (uh *UserHandlers) ListWithdraw(c *gin.Context) {
	list := requests.ListWithdraw(uh.DB, c)

	if list != nil {
		if len(list) > 0 {
			c.JSON(http.StatusOK, list)
			return
		}

		c.JSON(http.StatusNoContent, list)
		return
	}
}
