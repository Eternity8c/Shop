package client_cart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

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

func (c *Client) AddItem(ctx context.Context, userID, productID string, quantity int) error {
	const op = "cart_client.AddItem"

	body := map[string]int{"quantity": quantity}
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("%s: marshal body: %w", op, err)
	}

	url := fmt.Sprintf("%s/Cart/%s", c.baseURL, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("%s: unauthorized", op)
	case http.StatusBadRequest:
		// полезно прочитать тело — там может быть "out of stock" / "insufficient stock"
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s: bad request: %s", op, string(msg))
	default:
		return fmt.Errorf("%s: unexpected status %d", op, resp.StatusCode)
	}
}
