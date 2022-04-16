package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/genga911/yandex-praktikum-graduation-project/app/database/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const AuthCookieName = "auth"

type JWTAuth struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func Auth(db *database.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// проверим, есть ли у пользователя кука
		authCookieName := "auth"
		authCookie, err := c.Cookie(authCookieName)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			fmt.Println(fmt.Sprintf("Cookie error: %s", err))
			return
		}

		// проверим корректность куки
		// пример взят из официальной документации
		token, err := jwt.ParseWithClaims(authCookie, &JWTAuth{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Cookie token error: %v", token.Header["alg"])
			}

			return cfg.SecretKey, nil
		})

		var ID int
		if token != nil {
			if claims, ok := token.Claims.(*JWTAuth); ok && token.Valid {
				ID = claims.ID
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}

		// Убедимся что пользователь реально существует
		rp := repository.User{
			DB: db,
		}
		user, err := rp.Find(ID)
		if err != nil {
			c.AbortWithError(http.StatusForbidden, err)
			return
		}

		// запомним пользователя, чтобы не делать лищних запросов в БД
		c.Set("user", user)
	}
}

// SetAuthCookie установка авторизационной куки
func SetAuthCookie(user *models.User, c *gin.Context, cfg *config.Config) {
	claims := JWTAuth{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(cfg.AuthTTL))),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetCookie(AuthCookieName, tokenString, cfg.CookieTTL, "/", cfg.RunAddress, false, false)
}
