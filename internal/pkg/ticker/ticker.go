package ticker

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rajeshpudota/crypto-server/internal/pkg/config"
	"github.com/rajeshpudota/crypto-server/internal/pkg/errors"
)

const (
	Subscribe         string = "subscribe"
	TickerChannel1Sec string = "ticker/1s"
)

type TickerCache struct {
	sync.RWMutex
	Tickers map[string]*Ticker
}

func NewTickerCache() *TickerCache {
	return &TickerCache{
		Tickers: make(map[string]*Ticker),
	}
}

func (c *TickerCache) UpdateTickerCache(config config.Config, symbols []string) error {

	// Initialize HitBTC WebSocket connection
	ws, _, err := websocket.DefaultDialer.Dial(config.WebsocketBaseUrl, nil)
	if err != nil {
		log.Fatal("Failed to connect to HitBTC WebSocket:", err)
	}

	// Subscribe to ticker updates for the given symbols
	message := NewWebSocketRequest(Subscribe, TickerChannel1Sec, symbols)
	err = ws.WriteJSON(message)
	if err != nil {
		// TODO: remove using log, instead inject log object and use it
		log.Fatal("Failed to subscribe to ticker updates:", err)
	}

	// Start a goroutine to handle incoming WebSocket messages
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error reading WebSocket message:", err)
				continue
			}

			var data map[string]interface{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Println("error decoding WebSocket message:", err)
				continue
			}

			// Avoid acknowledgement response
			if _, ok := data["data"]; !ok {
				continue
			}

			// parse ticket response
			var ticketData WebSocketTicketResponse
			if err = json.Unmarshal(message, &ticketData); err != nil {
				log.Println("Error parsing WebSocket message:", err)
				continue
			}

			for symbol, data := range ticketData.Data {
				ticker := Ticker{
					Symbol: symbol,
					T:      data["t"].(float64),
					A:      data["a"].(string),
					B:      data["b"].(string),
					C:      data["c"].(string),
					O:      data["o"].(string),
					L:      data["l"].(string),
					H:      data["h"].(string),
				}
				c.UpdateTicker(ticker)
			}
			log.Println("Upated ticket cache")
		}
	}()

	return nil
}

func (c *TickerCache) GetTicker(symbol string) (*Ticker, error) {
	c.RLock()
	defer c.RUnlock()
	ticker, ok := c.Tickers[symbol]
	if !ok {
		return nil, errors.ErrTickerNotFound
	}
	return ticker, nil
}

func (c *TickerCache) GetAllTickers() []*Ticker {
	c.RLock()
	defer c.RUnlock()
	tickers := make([]*Ticker, 0, len(c.Tickers))
	for _, v := range c.Tickers {
		tickers = append(tickers, v)
	}
	return tickers
}

func (c *TickerCache) UpdateTicker(data Ticker) {
	c.Lock()
	defer c.Unlock()
	c.Tickers[data.Symbol] = &data
}
