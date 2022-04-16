package app

import (
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/handlers"
	"github.com/genga911/yandex-praktikum-graduation-project/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetUpServer запуск сервера
func SetUpServer() *gin.Engine {
	r := gin.Default()
	cfg := config.Get()
	config.InitFlags(cfg)

	db := database.GetDB(cfg)

	auth := middleware.GetAuthMiddleware(cfg)
	rApi := r.Group("/api")
	{
		uh := handlers.UserHandlers{DB: db}
		rUser := rApi.Group("/user")
		{
			// регистрация
			rUser.POST("/register", uh.Register)
			// авторизация
			rUser.POST("/login", auth.LoginHandler)

			// роуты под авторизацией
			rUser.Use(auth.MiddlewareFunc())
			{
				rUserOrders := rUser.Group("/orders")
				{
					// загрузка пользователем номера заказа для расчёта
					rUserOrders.POST("/")
					// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
					rUserOrders.GET("/")
				}

				rUserBalance := rUser.Group("/balance")
				{
					// получение текущего баланса счёта баллов лояльности пользователя
					rUserBalance.GET("/")
					// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
					rUserBalance.POST("/withdraw")
					// получение информации о выводе средств с накопительного счёта пользователем
					rUserBalance.GET("/withdrawals")
				}
			}
		}
	}

	err := r.Run(cfg.RunAddress)
	if err != nil {
		panic(fmt.Sprintf("Cannot run server: %s", err))
	}

	return r
}
