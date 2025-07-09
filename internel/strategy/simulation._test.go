package strategy

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestSimulation(t *testing.T) {

	positions := []Position{
		{
			Amount:    decimal.NewFromInt(25000),
			Leverage:  10,
			EntryDown: decimal.NewFromFloat(0.35),
		},
		{
			Amount:    decimal.NewFromInt(25000),
			Leverage:  10,
			EntryDown: decimal.NewFromFloat(0.4),
		},
		{
			Amount:    decimal.NewFromInt(25000),
			Leverage:  10,
			EntryDown: decimal.NewFromFloat(0.45),
		},
		{
			Amount:    decimal.NewFromInt(25000),
			Leverage:  10,
			EntryDown: decimal.NewFromFloat(0.5),
		},
	}
	basePrice := decimal.NewFromFloat(0.07976)
	currentPrice := decimal.NewFromFloat(0.042)
	pnl := Simulation(basePrice, currentPrice, positions)
	fmt.Println("总盈亏", pnl)
}
