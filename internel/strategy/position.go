package strategy

import "github.com/shopspring/decimal"

type Position struct {
	Amount    decimal.Decimal
	EntryDown decimal.Decimal
	Leverage  int64
}
