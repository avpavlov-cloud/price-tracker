package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"price-tracker/internal/models"
	"price-tracker/internal/repository"
	"price-tracker/internal/service"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
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

	// Инициализируем воркер
	worker := service.NewWorker(tracker)

	// Запускаем воркер в фоне
	go worker.Start(context.Background())

	ctx := context.Background()

	// Попробуем создать тестовый товар (просто для проверки)
	// В реальном приложении это будет идти через API
	fmt.Println("Обновляем цену для товара с ID 1...")

	err = tracker.TrackNewPrice(ctx, 1, 15500.50)
	if err != nil {
		log.Printf("Ошибка (возможно товара с ID 1 еще нет): %v", err)
	} else {
		fmt.Println("Цена успешно обновлена и сохранена в историю!")
	}

	e := echo.New()

	// Маршрут для создания товара
	e.POST("/products", func(c echo.Context) error {
		var p models.Product
		if err := c.Bind(&p); err != nil {
			return err
		}
		id, err := repo.CreateProduct(context.Background(), p)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		p.ID = id
		return c.JSON(http.StatusCreated, p)
	})

	// Маршрут для просмотра всех товаров
	e.GET("/products", func(c echo.Context) error {
		products, _ := repo.GetAllProducts(context.Background())
		return c.JSON(http.StatusOK, products)
	})

	// Добавил новый роут, чтобы смотреть историю цен конкретного товара
	e.GET("/products/:id/history", func(c echo.Context) error {
		// 1. Получаем ID из параметров пути
		idParam := c.Param("id")
		productID, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		}

		// 2. Запрашиваем историю у репозитория
		history, err := repo.GetPriceHistory(c.Request().Context(), productID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// 3. Возвращаем JSON массив
		return c.JSON(http.StatusOK, history)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
