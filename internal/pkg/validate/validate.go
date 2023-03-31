package validate

import "github.com/rajeshpudota/crypto-server/internal/pkg/symbol"

func IsValidSymbol(symbolCache *symbol.SymbolCache, symbol string) bool {
	if _, ok := symbolCache.Symbols[symbol]; ok {
		return true
	}
	return false
}
