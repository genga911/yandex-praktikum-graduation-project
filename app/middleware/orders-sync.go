package middleware

import (
	"errors"
	"net/http"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	"github.com/gin-gonic/gin"
)

func OrdersSync(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exist := c.Get("user")
		if !exist {
			c.AbortWithError(http.StatusUnauthorized, errors.New("пользователь не найден"))
			return
		}
		user := u.(*models.User)

		ro := repository.Order{
			DB: db,
		}

		err := ro.Sync(cfg, user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Next()
	}
}
