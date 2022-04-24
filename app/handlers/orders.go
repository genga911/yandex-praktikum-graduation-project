package handlers

import (
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/handlers/requests"
	"github.com/gin-gonic/gin"
)

type OrdersHandlers struct {
	DB  *database.DB
	Cfg *config.Config
}

func (oh *OrdersHandlers) OrdersUpload(c *gin.Context) {
	order := requests.OrderUpload(oh.DB, c)
	c.JSON(http.StatusAccepted, order)
}
