package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"price-tracker/internal/repository"
	"price-tracker/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := "postgres://user:password@localhost:5432/price_tracker?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	fmt.Println("Подключение к БД успешно, репозиторий готов!")

	tracker := service.NewTrackerService(repo)

	ctx := context.Background()

	// 1. Попробуем создать тестовый товар (просто для проверки)
	// В реальном приложении это будет идти через API
	fmt.Println("Обновляем цену для товара с ID 1...")

	err = tracker.TrackNewPrice(ctx, 1, 15500.50)
	if err != nil {
		log.Printf("Ошибка (возможно товара с ID 1 еще нет): %v", err)
	} else {
		fmt.Println("Цена успешно обновлена и сохранена в историю!")
	}
}
