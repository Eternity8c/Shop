package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Product struct {
	ID          int64  `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Price       int64  `json:"Price"`
	Stock       int64  `json:"Stock"`
}

type CreateProductDTO struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Price       int64  `json:"Price"`
	Stock       int64  `json:"Stock"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	log        *slog.Logger
}

func New(addr string, log *slog.Logger) *Client {
	return &Client{
		baseURL: "http://" + addr,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		log: log,
	}
}

func (c *Client) AllProducts(ctx context.Context) ([]Product, error) {
	const op = "Product.AllProducts"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/Products", nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var products []Product
	if err = json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (c *Client) ProductByName(ctx context.Context, name string) (Product, error) {
	const op = "Product.ProductByName"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/Product/search?name="+name, nil)
	if err != nil {
		return Product{}, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Product{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Product{}, fmt.Errorf("%s: product not found", op)
	}

	var product Product
	if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (c *Client) CreateProduct(ctx context.Context, dto CreateProductDTO) (int64, error) {
	const op = "Product.CreateProduct"

	b, err := json.Marshal(dto)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/created", bytes.NewReader(b))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var id int64
	if err = json.NewDecoder(resp.Body).Decode(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
