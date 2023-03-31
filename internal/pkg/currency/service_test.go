package currency

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
	service := NewCurrencyService(config, client)

	// Test case 1 - Successful request and response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"BTC": {
				"full_name": "Bitcoin TST",
				"crypto":true,
				"payin_enabled": true,
				"payout_enabled": true,
				"transfer_enabled": true,
				"precision_transfer": "0.00000001",
				"networks": [{
					"network": "BTC",
					"protocol": "OMNI",
					"default": true,
					"payin_enabled": true,
					"payout_enabled": true,
					"precision_payout": "0.00000001",
					"payout_fee": "0.000725840000",
					"payout_is_payment_id": false,
					"payin_payment_id": false,
					"payin_confirmations": 3
				}]
			},
			"ETH": {
				"full_name": "Ethereum TST",
				"crypto":true,
				"payin_enabled": true,
				"payout_enabled": true,
				"transfer_enabled": true,
				"precision_transfer": "0.000000000001",
				"networks": [{
					"network": "ETHTEST",
					"protocol": "",
					"default": true,
					"payin_enabled": true,
					"payout_enabled": true,
					"precision_payout": "0.000000000000000001",
					"payout_fee": "0.003621047265",
					"payout_is_payment_id": false,
					"payin_payment_id": false,
					"payin_confirmations": 2
				}]
			}
		}`))
	}))
	defer ts.Close()

	service.Config.APIBaseURL = ts.URL
	symbols, err := service.GetAll()
	if err != nil {
		t.Errorf("Test case 1 failed: expected nil error, but got %v", err)
	}
	if len(*symbols) != 2 {
		t.Errorf("Test case 1 failed: expected 2 symbols, but got %v", len(*symbols))
	}

	// Test case 2 - Request failure
	service.Config.APIBaseURL = "invalid_url"
	_, err = service.GetAll()
	if err == nil {
		t.Errorf("Test case 2 failed: expected non-nil error, but got nil")
	}
}
