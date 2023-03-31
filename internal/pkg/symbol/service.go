package symbol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rajeshpudota/crypto-server/internal/pkg/config"
)

type SymbolService struct {
	Config config.Config
	Client *http.Client
}

func NewSymbolService(config config.Config, client *http.Client) *SymbolService {
	return &SymbolService{config, client}
}

func (s *SymbolService) GetAll() (*[]Symbol, error) {
	url := fmt.Sprintf("%s/public/symbol", s.Config.APIBaseURL)

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

	var symbols map[string]Symbol
	err = json.Unmarshal(body, &symbols)
	if err != nil {
		return nil, err
	}

	var symbolList []Symbol
	for symbolId, symbol := range symbols {
		symbol.ID = symbolId
		symbolList = append(symbolList, symbol)
	}

	return &symbolList, nil
}
