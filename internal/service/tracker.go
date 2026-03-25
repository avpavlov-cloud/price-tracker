package service

import (
	"context"
	"price-tracker/internal/repository"

	"github.com/shopspring/decimal"
)

type TrackerService struct {
	repo *repository.Repository
}

func NewTrackerService(repo *repository.Repository) *TrackerService {
	return &TrackerService{repo: repo}
}

// TrackNewPrice — метод, который будет вызываться воркером или через API
func (s *TrackerService) TrackNewPrice(ctx context.Context, productID int, price float64) error {
	newPrice := decimal.NewFromFloat(price)

	// Здесь можно добавить проверку: если цена не изменилась, не делать лишний запрос в базу

	return s.repo.UpdatePriceTransaction(ctx, productID, newPrice)
}
