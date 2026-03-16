package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Product struct {
	ID          int64  `json:"ID"`
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

func (c *Client) GetProductByID(ctx context.Context, id int) (Product, error) {
	const op = "Client.GetProductByID"

	idStr := strconv.Itoa(id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/Products/"+idStr, nil)
	if err != nil {
		return Product{}, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Product{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var product Product
	if err = json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}
