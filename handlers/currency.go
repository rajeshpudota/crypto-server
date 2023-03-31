package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rajeshpudota/crypto-server/data"
	cache "github.com/rajeshpudota/crypto-server/internal/pkg/cache"
	"github.com/rajeshpudota/crypto-server/internal/pkg/errors"
	"github.com/rajeshpudota/crypto-server/internal/pkg/validate"
)

type Currency struct {
	l     *log.Logger
	cache *cache.Cache
}

// NewCurrency creates a currency handler with the given logger
func NewCurrency(l *log.Logger, cache *cache.Cache) *Currency {
	return &Currency{l, cache}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ListAll handles GET requests and returns all currencies currenlty
func (c *Currency) ListAll(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("[DEBUG] get all currencies")
	rw.Header().Add("Content-Type", "application/json")

	currencyList, err := c.cache.GetAllCurrencies()

	switch err {
	case nil:

	case errors.ErrSymbolNotFound:
		c.l.Println("[ERROR] fetching currency", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		c.l.Println("[ERROR] fetching curency", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(currencyList, rw)
	if err != nil {
		// we should never be here but log the error just incase
		c.l.Println("[ERROR] serializing product", err)
	}
}

// ListSingle handles GET requests
func (c *Currency) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	symbol := mux.Vars(r)["symbol"]

	if !validate.IsValidSymbol(c.cache.SymbolCache, symbol) {
		c.l.Printf("[ERROR] not a valid symbol(%s)", symbol)

		rw.WriteHeader(http.StatusBadGateway)
		data.ToJSON(&GenericError{Message: "symbol not valid"}, rw)
		return
	}

	currency, err := c.cache.GetCurrency(symbol)

	switch err {
	case nil:

	case errors.ErrSymbolNotFound:
		c.l.Println("[ERROR] fetching currency", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		c.l.Println("[ERROR] fetching curency", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(currency, rw)
	if err != nil {
		// we should never be here but log the error just incase
		c.l.Println("[ERROR] serializing product", err)
	}
}
