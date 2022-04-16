package handlers

import (
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/handlers/requests"
	"github.com/genga911/yandex-praktikum-graduation-project/app/middleware"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	DB  *database.DB
	Cfg *config.Config
}

func (uh *UserHandlers) Register(c *gin.Context) {
	// проверяем и создаем пользователя
	user := requests.DoRegister(uh.DB, c)
	if user != nil {
		// авторизуем его в системе
		middleware.SetAuthCookie(user, c, uh.Cfg)
		c.JSON(http.StatusOK, user)
	}
}
