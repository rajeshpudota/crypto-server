package errors

import "fmt"

var (
	ErrSymbolNotFound   = fmt.Errorf("symbol not found")
	ErrCurrencyNotFound = fmt.Errorf("currency not found")
	ErrTickerNotFound   = fmt.Errorf("ticker not found")
)
