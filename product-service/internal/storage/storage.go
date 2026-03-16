package storage

import (
	"context"
	"errors"
	"fmt"
	"product-service/internal/dto"
	"product-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "Storage.New"

	pool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		pool: pool,
	}, nil
}

func (s *Storage) ProductByName(ctx context.Context, name string) (models.Product, error) {
	const op = "Storage.ProductByName"

	stmt := `
	SELECT id, name, description, price, stock, created_at 
	FROM products
	WHERE name = $1
	`

	var product models.Product

	err := s.pool.QueryRow(ctx, stmt, name).
		Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, fmt.Errorf("%s: %w", op, ErrProductNotFound)
		}

		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (s *Storage) ProductByID(ctx context.Context, id int) (models.Product, error) {
	const op = "Storage.ProductByID"

	stmt := `
	SELECT id, name, description, price, stock, created_at 
	FROM products
	WHERE id = $1
	`

	var product models.Product

	err := s.pool.QueryRow(ctx, stmt, id).
		Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, fmt.Errorf("%s: %w", op, ErrProductNotFound)
		}

		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (s *Storage) AllProducts(ctx context.Context) ([]models.Product, error) {
	const op = "Storage.AllProducts"

	stmt := `
	SELECT id, name, description, price, stock, created_at 
	FROM products
	`

	rows, err := s.pool.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (s *Storage) CreatedProduct(ctx context.Context, dto dto.Product) (int64, error) {
	const op = "Storage.CreatedProduct"

	stmt := `
	INSERT INTO products (name, description, price, stock)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	var id int64

	err := s.pool.QueryRow(ctx, stmt, dto.Name, dto.Description, dto.Price, dto.Stock).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
