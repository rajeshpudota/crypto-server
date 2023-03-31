package ticker

type Ticker struct {
	Symbol string  `json:"symbol"`
	T      float64 `json:"t"` // Timestamp in milliseconds
	A      string  `json:"a"` // Best ask
	B      string  `json:"b"` // Best bid
	C      string  `json:"c"` // Best last price
	O      string  `json:"o"` // Open price
	L      string  `json:"l"` // Low price
	H      string  `json:"h"` // High price
}

type WebSocketRequest struct {
	Method  string `json:"method"`
	Channel string `json:"ch"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}

type Params struct {
	Symbols []string `json:"symbols"`
}

func NewWebSocketRequest(method, channel string, symbols []string) WebSocketRequest {
	return WebSocketRequest{
		Method:  method,
		Channel: channel,
		Params: Params{
			Symbols: symbols,
		},
	}
}

type WebSocketTicketResponse struct {
	Ch   string                            `json:"ch"`
	Data map[string]map[string]interface{} `json:"data"`
}
