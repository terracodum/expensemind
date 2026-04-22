package ml_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/terracodum/expensemind/backend/internal/ml"
)

func TestPredict_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/internal/v1/predict" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"forecast": []map[string]any{
				{"t": 3, "balance": 950.0},
				{"t": 4, "balance": 850.0},
			},
			"predicted_balance": 850.0,
			"confidence":        0.82,
		})
	}))
	defer srv.Close()

	client := ml.New(srv.URL)
	req := ml.PredictRequest{
		Timeseries: []ml.TimePoint{
			{T: 1, Balance: 1200.0},
			{T: 2, Balance: 1000.0},
		},
		Horizon:  2,
		Features: ml.Features{AvgDailyExpense: 180.0},
	}

	forecast, err := client.Predict(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(forecast.Points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(forecast.Points))
	}
	if forecast.Points[0].T != 3 || forecast.Points[0].Balance != 950.0 {
		t.Errorf("unexpected first point: %+v", forecast.Points[0])
	}
	if forecast.Points[1].T != 4 || forecast.Points[1].Balance != 850.0 {
		t.Errorf("unexpected second point: %+v", forecast.Points[1])
	}
	if forecast.PredictedBalance != 850.0 {
		t.Errorf("unexpected predicted_balance: %f", forecast.PredictedBalance)
	}
	if forecast.Confidence != 0.82 {
		t.Errorf("unexpected confidence: %f", forecast.Confidence)
	}
}

func TestPredict_serviceUnavailable(t *testing.T) {
	client := ml.New("http://127.0.0.1:19999") // ничего не слушает

	_, err := client.Predict(ml.PredictRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPredict_nonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"code": "PREDICTION_ERROR", "message": "internal error"},
		})
	}))
	defer srv.Close()

	client := ml.New(srv.URL)
	_, err := client.Predict(ml.PredictRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
