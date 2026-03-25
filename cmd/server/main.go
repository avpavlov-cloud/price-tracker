package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"price-tracker/internal/models"
	"price-tracker/internal/repository"
	"price-tracker/internal/service"
	"strconv"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
)

func main() {
	//  Создаем контекст, который отменится при Ctrl+C (SIGINT/SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
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
	go worker.Start(ctx)


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
		// Получаем ID из параметров пути
		idParam := c.Param("id")
		productID, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		}

		// Запрашиваем историю у репозитория
		history, err := repo.GetPriceHistory(c.Request().Context(), productID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Возвращаем JSON массив
		return c.JSON(http.StatusOK, history)
	})

	//Запускаем сервер в горутине, чтобы он не блокировал main
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Завершение работы сервера...")
		}
	}()

	// Программа замрет здесь, пока не придет сигнал
	<-ctx.Done() 
	fmt.Println("\nПолучен сигнал завершения. Плавная остановка...")

	//  Контекст с таймаутом (даем серверу 10 секунд на завершение дел)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Останавливаем Echo
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	//  Закрываем ресурсы (БД и воркеры)
	fmt.Println("Закрываем соединение с БД...")
	if err := db.Close(); err != nil {
		log.Printf("Ошибка при закрытии БД: %v", err)
	}

	fmt.Println("Сервер успешно остановлен.")
}
