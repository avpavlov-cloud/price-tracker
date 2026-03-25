package repository

import (
	"context"
	"database/sql"
	"price-tracker/internal/models"

	"github.com/shopspring/decimal"
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

// UpdatePriceTransaction обновляет цену товара и пишет историю в одной транзакции
func (r *Repository) UpdatePriceTransaction(ctx context.Context, productID int, newPrice decimal.Decimal) error {
	// Начинаем транзакцию
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Отложенный откат, если транзакция не завершится успехом (Commit)
	defer tx.Rollback()

	// 1. Обновляем текущую цену в таблице products
	queryUpdate := `UPDATE products SET current_price = $1 WHERE id = $2`
	if _, err := tx.ExecContext(ctx, queryUpdate, newPrice, productID); err != nil {
		return err
	}

	// 2. Добавляем запись в таблицу истории цен
	queryHistory := `INSERT INTO price_history (product_id, price) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, queryHistory, productID, newPrice); err != nil {
		return err
	}

	// Фиксируем изменения в базе
	return tx.Commit()
}

func (r *Repository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, url, current_price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.URL, &p.CurrentPrice); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetPriceHistory возвращает историю изменений цены для конкретного товара
func (r *Repository) GetPriceHistory(ctx context.Context, productID int) ([]models.PriceHistory, error) {
	query := `
		SELECT id, product_id, price, fetched_at 
		FROM price_history 
		WHERE product_id = $1 
		ORDER BY fetched_at DESC`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PriceHistory
	for rows.Next() {
		var h models.PriceHistory
		if err := rows.Scan(&h.ID, &h.ProductID, &h.Price, &h.FetchedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}


