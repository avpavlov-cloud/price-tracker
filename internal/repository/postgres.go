package repository

import (
	"context"
	"database/sql"
	"price-tracker/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateProduct добавляет новый товар в базу
func (r *Repository) CreateProduct(ctx context.Context, p models.Product) (int, error) {
	var id int
	query := `INSERT INTO products (name, url, current_price) 
              VALUES ($1, $2, $3) RETURNING id`
	
	err := r.db.QueryRowContext(ctx, query, p.Name, p.URL, p.CurrentPrice).Scan(&id)
	return id, err
}
