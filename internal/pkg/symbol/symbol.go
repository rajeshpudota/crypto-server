package symbol

import (
	"sync"

	"github.com/rajeshpudota/crypto-server/internal/pkg/errors"
)

type Symbol struct {
	ID            string `json:"id"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	FeeCurrency   string `json:"fee_currency"`
}

type SymbolCache struct {
	sync.RWMutex
	Symbols map[string]*Symbol
}

func NewSymbolCache() *SymbolCache {
	return &SymbolCache{}
}

func (s *SymbolCache) UpdaSymbolCacheCache(ss *SymbolService) error {
	s.Symbols = map[string]*Symbol{}
	symbols, err := ss.GetAll()

	if err != nil {
		return err
	}

	for _, symbol := range *symbols {
		s.updateSymbol(symbol)
	}

	return nil
}

func (s *SymbolCache) GetSymbol(symbol string) (*Symbol, error) {
	s.RLock()
	defer s.RUnlock()
	symbolData, ok := s.Symbols[symbol]
	if !ok {
		return nil, errors.ErrSymbolNotFound
	}
	return symbolData, nil
}

func (s *SymbolCache) updateSymbol(data Symbol) {
	s.Lock()
	defer s.Unlock()
	s.Symbols[data.ID] = &data
}
