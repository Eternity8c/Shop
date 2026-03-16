package service

import (
	"cart-service/internal/client"
	"cart-service/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"strconv"
)

type Cart struct {
	log     *slog.Logger
	storage CartStorage
	client  ProductClient
}

type CartStorage interface {
	AddItem(ctx context.Context, userID string, item domain.CartItem) error
}

type ProductClient interface {
	GetById(ctx context.Context, id int) (client.Product, error)
}

func (c *Cart) Add(ctx context.Context, userID string, productID int) (bool, error) {
	const op = "Cart.Add"

	product, err := c.client.GetById(ctx, productID)
	if err != nil {
		return false, fmt.Errorf("%s: failed to get product: %w", op, err)
	}

	if product.Stock <= 0 {
		return false, fmt.Errorf("%s: product out of stock", op)
	}

	cartItem := domain.CartItem{
		ProductID: strconv.Itoa(productID),
		Name:      product.Name,
		Price:     float64(product.Price),
		Quantity:  1,
	}

	if err := c.storage.AddItem(ctx, userID, cartItem); err != nil {
		return false, fmt.Errorf("%s: failed to add item to cart: %w", op, err)
	}

	return true, nil
}
