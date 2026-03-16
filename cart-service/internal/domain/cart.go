package domain

import (
	"encoding/json"
	"time"
)

type CartItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	ImageURL  string  `json:"image_url"`
}

type Cart struct {
	UserID     string     `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
	TotalItems int        `json:"total_items"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CartSummary для быстрого получения информации
type CartSummary struct {
	UserID     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	TotalItems int       `json:"total_items"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Cart) CalculateTotals() {
	c.TotalPrice = 0
	c.TotalItems = 0
	for _, item := range c.Items {
		c.TotalPrice += item.Price * float64(item.Quantity)
		c.TotalItems += item.Quantity
	}
}

func (c *Cart) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Cart) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}
