package requests

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	request_errors "github.com/genga911/yandex-praktikum-graduation-project/app/handlers/requests/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/theplant/luhn"
)

func OrderUpload(db *database.DB, cfg *config.Config, c *gin.Context) *models.Order {
	body := c.Request.Body
	number, err := ioutil.ReadAll(body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	numberString := string(number)
	// Если валидация провалилась, просто выйдем
	ok := validateNumber(numberString, c)
	if !ok {
		return nil
	}

	u, exist := c.Get("user")
	if !exist {
		c.AbortWithError(http.StatusUnauthorized, errors.New("пользователь не найден"))
		return nil
	}

	user := u.(*models.User)
	rp := repository.Order{
		DB: db,
	}

	order := models.Order{Number: numberString}
	err = rp.Create(user, &order)

	if err != nil {
		var pgErr *request_errors.UniqError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				if pgErr.Message == request_errors.OrderAlreadyExists {
					c.AbortWithError(http.StatusOK, errors.New(request_errors.OrderAlreadyExists))
				} else {
					c.AbortWithError(http.StatusConflict, errors.New(request_errors.OrderCreatedByAnotherUser))
				}
				return &order
			}
		}

		// если не попали под прошлые условия выбросим общую ошибку
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}

	return &order
}

// OrdersList список заказов
func OrdersList(db *database.DB, c *gin.Context) []*models.Order {
	rp := repository.Order{
		DB: db,
	}

	u, exist := c.Get("user")
	if !exist {
		c.AbortWithError(http.StatusUnauthorized, errors.New("пользователь не найден"))
		return nil
	}
	user := u.(*models.User)
	orders, err := rp.List(user, true)

	if err != nil {
		return []*models.Order{}
	}

	return orders
}

// validateNumber Валидатор для строки с номером заказа
func validateNumber(number string, c *gin.Context) bool {
	if len(number) == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("пустой номер заказа"))
		return false
	}

	matched, err := regexp.MatchString(`^\d+$`, number)
	if err != nil || !matched {
		c.AbortWithError(http.StatusUnprocessableEntity, errors.New("номер должен состоять только из цифр"))
		return false
	}

	orderNumber, _ := strconv.Atoi(number)
	if !luhn.Valid(orderNumber) {
		c.AbortWithError(http.StatusUnprocessableEntity, errors.New("номер не корректен"))
		return false
	}

	return true
}
