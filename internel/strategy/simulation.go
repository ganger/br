package strategy

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func Simulation(basePrice, currentPrice decimal.Decimal, positions []Position) (totalPnl decimal.Decimal) {
	profit := decimal.Zero
	loss := decimal.Zero
	for _, position := range positions {
		entryPrice := basePrice.Mul(decimal.NewFromInt(1).Sub(position.EntryDown))
		fmt.Println("======")
		fmt.Println("挂单价格:", entryPrice)
		if entryPrice.LessThan(currentPrice) {
			fmt.Println("未达价格，未成交")
			continue
		}
		liquidationPrice, _ := CalculateLiquidationPrice(entryPrice, position.Amount, position.Leverage)
		if liquidationPrice.GreaterThan(currentPrice) {
			loss = loss.Add(position.Amount.Div(decimal.NewFromInt(position.Leverage)))
			fmt.Println("爆仓，总损失:", loss, " 爆仓价格:", liquidationPrice)
			continue
		}
		profit = position.Amount.Mul(basePrice.Sub(entryPrice).Div(entryPrice))
		fmt.Println("挣钱:", profit, "开仓数量:", position.Amount.Div(entryPrice))
	}
	totalPnl = profit.Sub(loss)
	return
}
