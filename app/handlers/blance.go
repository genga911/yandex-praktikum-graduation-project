package handlers

import (
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/handlers/requests"
	"github.com/gin-gonic/gin"
)

// Balance Баланс пользователя
func (uh *UserHandlers) Balance(c *gin.Context) {
	balance := requests.GetBalance(uh.Cfg, uh.DB, c)
	c.JSON(http.StatusOK, balance)
}

func (uh *UserHandlers) RegisterWithdraw(c *gin.Context) {
	withdraw := requests.RegisterWithdraw(uh.Cfg, uh.DB, c)

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
