package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"product-service/internal/dto"
	"product-service/internal/models"
	"product-service/internal/storage"
)

type Service struct {
	log         *slog.Logger
	productStor ProductStorage
}

type ProductStorage interface {
	CreatedProduct(ctx context.Context, dto dto.Product) (int64, error)
	AllProducts(ctx context.Context) ([]models.Product, error)
	ProductByName(ctx context.Context, name string) (models.Product, error)
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductData     = errors.New("invalid product data")
)

func New(log *slog.Logger, pr ProductStorage) *Service {
	return &Service{
		productStor: pr,
		log:         log,
	}
}

func (s *Service) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	const op = "Service.GetAllProducts"

	products, err := s.productStor.AllProducts(ctx)
	if err != nil {
		s.log.Warn("%s: %w", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (s *Service) GetProductByName(ctx context.Context, name string) (models.Product, error) {
	const op = "Service.GetProductByName"

	if err := validatorName(name); err != nil {
		s.log.Error("%s: %w", op, err)
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	product, err := s.productStor.ProductByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return models.Product{}, fmt.Errorf("%s: %w", op, ErrProductNotFound)
		}
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (s *Service) CreatedProduct(ctx context.Context, dto dto.Product) (int64, error) {
	const op = "Service.CreatedProduct"

	if err := validatorDTO(dto); err != nil {
		s.log.Error("%s: %w", op, err)
		return 0, fmt.Errorf("%s: %w", op, ErrProductData)
	}

	id, err := s.productStor.CreatedProduct(ctx, dto)
	if err != nil {
		s.log.Error("%s: %w", op, err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func validatorName(name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}

func validatorDTO(dto dto.Product) error {
	if dto.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if dto.Description == "" {
		return fmt.Errorf("Description is empty")
	}
	if dto.Price == 0 {
		return fmt.Errorf("Price is empty")
	}
	if dto.Stock == 0 {
		return fmt.Errorf("Stock is empty")
	}
	return nil
}
