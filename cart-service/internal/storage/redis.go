package storage

import (
	"cart-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

const (
	cartKeyPrefix     = "cart:"         // String: полные данные корзины
	cartSummaryPrefix = "cart:summary:" // Hash: краткая информация
	cartExpireTime    = 24 * time.Hour  // Время жизни корзины
	cartLockPrefix    = "cart:lock:"    // Для блокировок
)

func NewStorage(ctx context.Context, storagePath string, log *slog.Logger) (*Storage, error) {
	const op = "Storage.NewStorage"

	opt, err := redis.ParseURL(storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		client: client,
	}, err
}

func (s *Storage) Save(ctx context.Context, cart *domain.Cart) error {
	const op = "Storage.Save"

	key := fmt.Sprintf("%s%s", cartKeyPrefix, cart.UserID)

	cart.UpdatedAt = time.Now()
	cart.CalculateTotals()

	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.client.Set(ctx, key, data, cartExpireTime).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	summaryKey := fmt.Sprintf("%s%s", cartSummaryPrefix, cart.UserID)
	summary := map[string]interface{}{
		"total_price": cart.TotalPrice,
		"total_items": cart.TotalItems,
		"update_at":   cart.UpdatedAt,
	}

	if err := s.client.HSet(ctx, summaryKey, summary).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.client.Expire(ctx, summaryKey, cartExpireTime)

	return nil
}

func (s *Storage) AddItem(ctx context.Context, userID string, item domain.CartItem) error {
	const op = "Storage.AddItem"

	cartKey := fmt.Sprintf("cart:%s", userID)

	pipe := s.client.Pipeline()

	data, err := s.client.Get(ctx, cartKey).Bytes()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var cart domain.Cart

	if err == redis.Nil {
		cart = domain.Cart{
			UserID:    userID,
			Items:     []domain.CartItem{},
			UpdatedAt: time.Now(),
		}
	} else {
		if err := json.Unmarshal(data, &cart); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, item)
	}

	cart.CalculateTotals()
	cart.UpdatedAt = time.Now()

	updatedData, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("error marshaling cart: %w", err)
	}

	pipe.Set(ctx, cartKey, updatedData, 24*time.Hour)

	summaryKey := fmt.Sprintf("cart:summary:%s", userID)
	pipe.HSet(ctx, summaryKey, map[string]interface{}{
		"total_price": cart.TotalPrice,
		"total_items": cart.TotalItems,
		"updated_at":  cart.UpdatedAt.Unix(),
	})
	pipe.Expire(ctx, summaryKey, 24*time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("error executing pipeline: %w", err)
	}

	return nil
}
