package symbol

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rajeshpudota/crypto-server/internal/pkg/config"
)

func TestGetAll(t *testing.T) {
	// Setup
	config := config.Config{APIBaseURL: "http://localhost:8080"}
	client := &http.Client{}
	service := NewSymbolService(config, client)

	// Test case 1 - Successful request and response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"ETHBTC": {
				"type": "spot",
				"base_currency": "ETH",
				"quote_currency": "BTC",
				"status": "working",
				"quantity_increment": "0.001",
				"tick_size": "0.000001",
				"take_rate": "0.001",
				"make_rate": "-0.0001",
				"fee_currency": "BTC",
				"margin_trading": true,
				"max_initial_leverage": "10.00"
			}
		}`))
	}))
	defer ts.Close()

	service.Config.APIBaseURL = ts.URL
	symbols, err := service.GetAll()
	if err != nil {
		t.Errorf("Test case 1 failed: expected nil error, but got %v", err)
	}
	if len(*symbols) != 1 {
		t.Errorf("Test case 1 failed: expected 1 symbols, but got %v", len(*symbols))
	}

	// Test case 2 - Request failure
	service.Config.APIBaseURL = "invalid_url"
	_, err = service.GetAll()
	if err == nil {
		t.Errorf("Test case 2 failed: expected non-nil error, but got nil")
	}
}
