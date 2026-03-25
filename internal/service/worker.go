package service

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type Worker struct {
	tracker *TrackerService
}

func NewWorker(ts *TrackerService) *Worker {
	return &Worker{tracker: ts}
}

// Start запускает бесконечный цикл обновления цен
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Интервал обновления
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	log.Println("Воркер: Начинаю обновление цен...")

	// Получаем все товары через репозиторий (нужно добавить метод в репозиторий)
	products, err := w.tracker.repo.GetAllProducts(ctx)
	if err != nil {
		log.Printf("Воркер: ошибка получения товаров: %v", err)
		return
	}

	for _, p := range products {
		// Имитируем изменение цены: +/- 5%
		change := p.CurrentPrice.InexactFloat64() * (rand.Float64()*0.1 - 0.05)
		newPrice := p.CurrentPrice.InexactFloat64() + change

		log.Printf("Воркер: Обновляю товар %s. Новая цена: %.2f", p.Name, newPrice)

		err := w.tracker.TrackNewPrice(ctx, p.ID, newPrice)
		if err != nil {
			log.Printf("Воркер: ошибка сохранения цены для ID %d: %v", p.ID, err)
		}
	}
}
