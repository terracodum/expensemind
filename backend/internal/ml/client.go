package ml

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/terracodum/expensemind/backend/internal/domain"
	"github.com/terracodum/expensemind/backend/internal/errors"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{},
	}
}

func (c *Client) Predict(req PredictRequest) (domain.Forecast, error) {
	body, _ := json.Marshal(req)

	resp, err := c.http.Post(c.baseURL+"/internal/v1/predict", "application/json", bytes.NewReader(body))
	if err != nil {
		return domain.Forecast{}, errors.MLServiceUnavailable("ml service unavailable", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Forecast{}, errors.MLResponseInvalid("unexpected status from ml service", nil)
	}

	var r predictResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return domain.Forecast{}, errors.MLResponseInvalid("failed to decode ml response", err)
	}

	points := make([]domain.Point, len(r.Forecast))
	for i, p := range r.Forecast {
		points[i] = domain.Point{T: p.T, Balance: p.Balance}
	}

	return domain.Forecast{
		Points:           points,
		PredictedBalance: r.PredictedBalance,
		Confidence:       r.Confidence,
	}, nil
}
