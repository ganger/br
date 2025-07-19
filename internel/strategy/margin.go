package strategy

import (
	"github.com/shopspring/decimal"
)

type MarginInfo struct {
	MaintenanceRate      decimal.Decimal
	QuickCalculateAmount decimal.Decimal
}

var MarginMap = map[int64]MarginInfo{
	20: {
		MaintenanceRate:      decimal.NewFromFloat(0.025),
		QuickCalculateAmount: decimal.NewFromFloat(27.5),
	},
	10: {
		MaintenanceRate:      decimal.NewFromFloat(0.05),
		QuickCalculateAmount: decimal.NewFromFloat(277.5),
	},
}

func CalculateMargin(amount decimal.Decimal, leverage int64) (initialMargin, maintenanceMargin decimal.Decimal) {
	info, ok := MarginMap[leverage]
	if !ok {
		return decimal.Zero, decimal.Zero
	}
	initialMargin = amount.Div(decimal.NewFromInt(leverage))
	maintenanceMargin = amount.Mul(info.MaintenanceRate).Sub(info.QuickCalculateAmount)
	return
}

func CalculateLiquidationPrice(entryPrice, amount decimal.Decimal, leverage int64) (liquidationPriceDown, liquidationPriceUp, down decimal.Decimal) {
	initialMargin, maintenanceMargin := CalculateMargin(amount, leverage)
	loss := initialMargin.Sub(maintenanceMargin)
	down = loss.Div(amount)
	liquidationPriceDown = entryPrice.Mul(decimal.NewFromInt(1).Sub(down))
	liquidationPriceUp = entryPrice.Mul(decimal.NewFromInt(1).Add(down))
	return
}
