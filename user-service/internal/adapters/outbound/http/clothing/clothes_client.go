package httpclothing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"user-service/internal/application/ports"
	"user-service/internal/domain"
)

var _ ports.ClothesClient = (*Client)(nil)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetByID(ctx context.Context, id int64) (*domain.Clothes, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/clothes/%d", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("clothing service returned status %d", resp.StatusCode)
	}

	var apiResp clothesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("clothing service returned unsuccessful response")
	}

	return &domain.Clothes{
		ID:      apiResp.Payload.ID,
		Name:    apiResp.Payload.Name,
		Price:   apiResp.Payload.Price,
		ColorID: apiResp.Payload.ColorID,
		Images:  toDomainImages(apiResp.Payload.Images),
	}, nil
}

type clothesAPIResponse struct {
	Success bool                 `json:"success"`
	Payload clothesPayload       `json:"payload"`
	Error   *clothesErrorPayload `json:"error,omitempty"`
}

type clothesErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type clothesPayload struct {
	ID      int64                 `json:"id"`
	Name    string                `json:"name"`
	Price   float64               `json:"price"`
	ColorID *int64                `json:"color_id"`
	Images  []clothesImagePayload `json:"images"`
}

type clothesImagePayload struct {
	ID        int64  `json:"id"`
	ClothesID int64  `json:"clothes_id"`
	ImageURL  string `json:"image_url"`
}

func toDomainImages(images []clothesImagePayload) []domain.ClothesImage {
	result := make([]domain.ClothesImage, 0, len(images))
	for _, image := range images {
		result = append(result, domain.ClothesImage{
			ID:        image.ID,
			ClothesID: image.ClothesID,
			ImageURL:  image.ImageURL,
		})
	}

	return result
}
