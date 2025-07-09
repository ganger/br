package strategy

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestGetTokenPrice(t *testing.T) {
	amount := int64(25000)

	initialMargin, maintenanceMargin := CalculateMargin(decimal.NewFromInt(amount), 10)

	fmt.Println("initialMargin:", initialMargin)
	fmt.Println("maintenanceMargin:", maintenanceMargin)
	loss := initialMargin.Sub(maintenanceMargin)
	fmt.Println("loss:", loss)
	down := loss.Div(decimal.NewFromInt(amount))
	fmt.Println("down:", down)

}
func TestCalculateLiquidationPrice(t *testing.T) {
	amount := int64(25000)
	entryPrice := decimal.NewFromFloat(0.1)
	liquidationPrice, down := CalculateLiquidationPrice(entryPrice, decimal.NewFromInt(amount), 10)
	fmt.Println("liquidationPrice:", liquidationPrice)
	fmt.Println("down:", down)
}
