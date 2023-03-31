package cache

import (
	"log"
	"net/http"
	"time"

	"github.com/rajeshpudota/crypto-server/internal/pkg/config"
	"github.com/rajeshpudota/crypto-server/internal/pkg/currency"
	"github.com/rajeshpudota/crypto-server/internal/pkg/symbol"
	"github.com/rajeshpudota/crypto-server/internal/pkg/ticker"
)

type Cache struct {
	CurrencyCache *currency.CurrencyCache
	TickerCache   *ticker.TickerCache
	SymbolCache   *symbol.SymbolCache
	l             *log.Logger
}

func NewCache(l *log.Logger) *Cache {
	return &Cache{
		CurrencyCache: currency.NewCurrencyCache(),
		TickerCache:   ticker.NewTickerCache(),
		SymbolCache:   symbol.NewSymbolCache(),
		l:             l,
	}
}

func (c *Cache) UpdateCache(config config.Config, symbols []string) error {
	// TODO: handle errors
	// Http Client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Currency: updates the cache every hour
	currencyService := currency.NewCurrencyService(config, client)
	go func() {
		for {
			err := c.CurrencyCache.UpdaCurrencyCacheCache(currencyService)
			c.l.Println("[Info] Currency caching complete")
			if err != nil {
				c.l.Printf("Error updating currency cache: %v", err)
			}
			time.Sleep(time.Hour)
		}
	}()

	// Symbols: updates the cache every hour
	symbolService := symbol.NewSymbolService(config, client)
	go func() {
		for {
			err := c.SymbolCache.UpdaSymbolCacheCache(symbolService)
			c.l.Println("[Info] Symbol caching complete")
			if err != nil {
				c.l.Printf("Error updating symbol cache: %v", err)
			}
			time.Sleep(time.Hour)
		}
	}()

	// Tickers
	c.TickerCache.UpdateTickerCache(config, symbols)
	c.l.Println("[Info] Intial ticker caching done for symbols: ", symbols)
	return nil
}

func (c *Cache) GetAllCurrencies() (*[]CurrencyResponse, error) {
	// Get all tickers
	tickers := c.TickerCache.GetAllTickers()
	currencies := []CurrencyResponse{}

	// Generate Currency in required response format
	for _, ticker := range tickers {
		// Get symbol
		symbol, err := c.SymbolCache.GetSymbol(ticker.Symbol)
		if err != nil {
			return nil, err
		}

		// Get Currency
		currency, err := c.CurrencyCache.GetCurrency(symbol.BaseCurrency)
		if err != nil {
			return nil, err
		}

		response := CurrencyResponse{
			ID:          currency.ID,
			FullName:    currency.FullName,
			Ask:         ticker.A,
			Bid:         ticker.B,
			Last:        ticker.C,
			Open:        ticker.O,
			Low:         ticker.L,
			High:        ticker.H,
			FeeCurrency: symbol.FeeCurrency,
		}
		currencies = append(currencies, response)
	}
	return &currencies, nil
}

func (c *Cache) GetCurrency(symbol string) (*CurrencyResponse, error) {
	// Get Ticker
	ticker, err := c.TickerCache.GetTicker(symbol)
	if err != nil {
		return nil, err
	}

	// Get symbol
	symbolObject, err := c.SymbolCache.GetSymbol(symbol)
	if err != nil {
		return nil, err
	}

	// Get Currency
	currency, err := c.CurrencyCache.GetCurrency(symbolObject.BaseCurrency)
	if err != nil {
		return nil, err
	}
	return &CurrencyResponse{
		ID:          currency.ID,
		FullName:    currency.FullName,
		Ask:         ticker.A,
		Bid:         ticker.B,
		Last:        ticker.C,
		Open:        ticker.O,
		Low:         ticker.L,
		High:        ticker.H,
		FeeCurrency: symbolObject.FeeCurrency,
	}, nil
}
