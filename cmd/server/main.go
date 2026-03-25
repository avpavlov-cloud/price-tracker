package main

import (
	"database/sql"
	"fmt"
	"log"
	"price-tracker/internal/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := "postgres://user:password@localhost:5432/price_tracker?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	repository.NewRepository(db)
	fmt.Println("Подключение к БД успешно, репозиторий готов!")

	// Дальше здесь будет запуск HTTP сервера или воркера
}
