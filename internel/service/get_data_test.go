package service

import (
	"br-trade/bootstrap"
	"br-trade/constx"
	"br-trade/global"
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
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
	s.AvgBrPrice = decimal.NewFromFloat(0.1)
	result := s.GetPriceToAvgSpreadPct()
	fmt.Println(result)
}

func TestGetBrPoolBalanceLow(t *testing.T) {
	s := NewDataService()

	s.PoolInfo.BrBalance = decimal.NewFromInt(3001)
	s.PoolInfo.AvgBrBalance = decimal.NewFromInt(10000)
	result := s.GetBrPoolBalanceLow()
	fmt.Println(result)
}

func TestGetUsdtPoolBalanceLow(t *testing.T) {
	s := NewDataService()

	s.PoolInfo.UsdtBalance = decimal.NewFromInt(3000)
	s.PoolInfo.AvgUsdtBalance = decimal.NewFromInt(10000)
	result := s.GetUsdtPoolBalanceLow()
	fmt.Println(result)
}

func TestCreateOrder(t *testing.T) {

	bootstrap.InitLogger()
	s := NewDataService()
	s.AvgBrPrice = decimal.NewFromFloat(0.0716)
	s.CreateOrder(futures.SideTypeBuy)
	fmt.Println("=============")
	s.CreateOrder(futures.SideTypeSell)
}

func TestBinanceOrder(t *testing.T) {
	bootstrap.InitConfig()
	bootstrap.InitBscClient()
	bootstrap.InitLogger()
	_, err := global.BinanceFuturesClient.NewCreateOrderService().
		Symbol(constx.BrFutureSymbol).
		Side(futures.SideTypeBuy).
		Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).
		Price("0.001").
		Quantity("10000").
		Do(context.Background())
	if err != nil {
		global.Logger.Error(err.Error())
	}
	fmt.Println("success")
}
