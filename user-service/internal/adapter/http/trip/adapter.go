package trip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/korroziea/taxi/user-service/internal/config"
	"github.com/korroziea/taxi/user-service/internal/domain"
)

const (
	requestTimeout = 10 * time.Second

	getTripsPath = "/api/trips"
)

type Adapter struct {
	baseURL string
	client  *http.Client
}

func New(cfg config.Trip) *Adapter {
	return &Adapter{
		baseURL: cfg.URL,
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (a *Adapter) Trips(ctx context.Context, userID string) ([]domain.Trip, error) {
	reqData := toTripsReq(userID)

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return []domain.Trip{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf(a.baseURL + getTripsPath)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return []domain.Trip{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return []domain.Trip{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return []domain.Trip{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	var respData []tripResp
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return []domain.Trip{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return toDomains(respData), nil
}
