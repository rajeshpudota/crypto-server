package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rajeshpudota/crypto-server/internal/pkg/config"
)

type CurrencyService struct {
	Config config.Config
	Client *http.Client
}

func NewCurrencyService(config config.Config, client *http.Client) *CurrencyService {
	return &CurrencyService{config, client}
}

func (s *CurrencyService) GetAll() (*[]Currency, error) {
	url := fmt.Sprintf("%s/public/currency", s.Config.APIBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d status code", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var currencies map[string]Currency
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return nil, err
	}

	var currencyList []Currency
	for currencyId, currency := range currencies {
		currencyList = append(currencyList, Currency{ID: currencyId, FullName: currency.FullName})
	}

	return &currencyList, nil
}
