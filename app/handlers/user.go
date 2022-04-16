package handlers

import (
	"github.com/genga911/yandex-praktikum-graduation-project/app/database"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	DB *database.DB
}

func (uh *UserHandlers) Register(c *gin.Context) {

}
