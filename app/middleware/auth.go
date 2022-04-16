package middleware

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/genga911/yandex-praktikum-graduation-project/app/config"
	models2 "github.com/genga911/yandex-praktikum-graduation-project/app/database/models"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// нагрузка токена
func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models2.User); ok {
		return jwt.MapClaims{
			identityKey: v.Name,
		}
	}
	return jwt.MapClaims{}
}

// аутентификация
func authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &models2.User{
			Name: userID,
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

// идентификация, получение данных из токена
func identity(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models2.User{
		Name: claims[identityKey].(string),
	}
}

// авторизация
func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*models2.User); ok && v.Name == "admin" {
		return true
	}

	return false
}

// разавторизация
func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func GetAuthMiddleware(cfg *config.Config) *jwt.GinJWTMiddleware {
	m, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:             []byte(cfg.JWTKey),
		Timeout:         time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     payload,
		IdentityHandler: identity,
		Authenticator:   authenticator,
		Authorizator:    authorizator,
		Unauthorized:    unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		panic(fmt.Sprintf("JWT Error: %s", err))
	}

	errInit := m.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return m
}
