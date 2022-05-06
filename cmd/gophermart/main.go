package main

import (
	"fmt"

	"github.com/genga911/yandex-praktikum-graduation-project/app"
	"github.com/joho/godotenv"
)

func main() {
	// загруаем данные из env
	if err := godotenv.Load(); err != nil {
		// если нет файла, просто пропустим. не страшно
		fmt.Printf("Load .env error: %s", err)
	}

	app.SetUpServer()
}
