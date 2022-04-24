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

	rAPI := r.Group("/api")
	{
		uh := handlers.UserHandlers{DB: db, Cfg: cfg}
		oh := handlers.OrdersHandlers{DB: db, Cfg: cfg}

		rUser := rAPI.Group("/user")
		{
			// регистрация
			rUser.POST("/register", uh.Register)
			// авторизация
			rUser.POST("/login", uh.Login)

			// роуты под авторизацией
			rUser.Use(middleware.Auth(db, cfg))
			{
				rUserOrders := rUser.Group("/orders")
				{
					// загрузка пользователем номера заказа для расчёта
					rUserOrders.POST("/", oh.OrdersUpload)
					// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
					rUserOrders.GET("/", oh.OrdersList)
				}

				rUserBalance := rUser.Group("/balance")
				{
					// получение текущего баланса счёта баллов лояльности пользователя
					rUserBalance.GET("/", uh.Balance)
					// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
					rUserBalance.POST("/withdraw", uh.RegisterWithdraw)
					// получение информации о выводе средств с накопительного счёта пользователем
					rUserBalance.GET("/withdrawals", uh.ListWithdraw)
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
