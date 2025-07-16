package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestGetBasisPct(t *testing.T) {
	s := NewDataService()
	s.BrPrice = decimal.NewFromFloat(0.1)
	s.BrFuturePrice = decimal.NewFromFloat(0.11)
	result := s.GetBasisPct()
	fmt.Println(result)
}

func TestGetPriceToAvgSpreadPct(t *testing.T) {
	s := NewDataService()

	s.BrPrice = decimal.NewFromFloat(0.11)
	avgPrice := decimal.NewFromFloat(0.1)
	result := s.getPriceToAvgSpreadPct(avgPrice)
	fmt.Println(result)
}

func TestGetBrPoolBalanceLow(t *testing.T) {
	s := NewDataService()

	s.PoolInfo.BrBalance = decimal.NewFromInt(2000)
	avgBalance := decimal.NewFromInt(10000)
	result := s.getBrPoolBalanceLow(avgBalance)
	fmt.Println(result)
}

func TestGetUsdtPoolBalanceLow(t *testing.T) {
	s := NewDataService()

	s.PoolInfo.UsdtBalance = decimal.NewFromInt(3001)
	avgBalance := decimal.NewFromInt(10000)
	result := s.getUsdtPoolBalanceLow(avgBalance)
	fmt.Println(result)
}
