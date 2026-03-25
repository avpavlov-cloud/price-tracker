package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID           int             `db:"id" json:"id"`
	Name         string          `db:"name" json:"name"`
	URL          string          `db:"url" json:"url"`
	CurrentPrice decimal.Decimal `db:"current_price" json:"current_price"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
}
