package currency

import (
	"sync"

	"github.com/rajeshpudota/crypto-server/internal/pkg/errors"
)

type Currency struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type CurrencyCache struct {
	sync.RWMutex
	Currency map[string]*Currency
}

func NewCurrencyCache() *CurrencyCache {
	return &CurrencyCache{}
}

func (c *CurrencyCache) UpdaCurrencyCache(s *CurrencyService) error {
	c.Currency = map[string]*Currency{}
	currencies, err := s.GetAll()
	if err != nil {
		return err
	}

	for _, currency := range *currencies {
		c.updateCurrency(currency)
	}
	return nil
}

func (c *CurrencyCache) GetCurrency(id string) (*Currency, error) {
	c.RLock()
	defer c.RUnlock()
	currencyData, ok := c.Currency[id]
	if !ok {
		return nil, errors.ErrCurrencyNotFound
	}
	return currencyData, nil
}

func (c *CurrencyCache) updateCurrency(data Currency) {
	c.Lock()
	defer c.Unlock()
	c.Currency[data.ID] = &data
}
