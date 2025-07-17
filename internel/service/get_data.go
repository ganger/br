package service

import (
	"br-trade/constx"
	"br-trade/global"
	"br-trade/internel/data"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

func (s *DataService) RefreshBrPrice() {
	price, err := data.GetTokenPriceFromBinance(constx.BrAddress)
	if err != nil {
		global.Logger.Error("获取现货价格异常")
		global.Logger.Error(err.Error())
		return
	}
	s.BrPrice = price
	s.SaveData(constx.RedisKeyBrPrice, fmt.Sprintf("%s", s.BrPrice.String()), time.Now())

	avgPrice, err := s.GetData(constx.RedisKeyBrPrice)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}
	s.AvgBrPrice = avgPrice
}

func (s *DataService) RefreshBrFuturePrice() {
	ts, buy, _, err := data.GetFuturePrice(constx.BrFutureSymbol)
	if err != nil {
		global.Logger.Error("获取期货价格异常")
		global.Logger.Error(err.Error())
		return
	}
	s.BrFuturePrice = buy

	t := time.UnixMilli(ts)
	s.SaveData(constx.RedisKeyBrFuturePrice, fmt.Sprintf("%s", s.BrFuturePrice.String()), t)

	avgPrice, err := s.GetData(constx.RedisKeyBrFuturePrice)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}
	s.AvgBrFuturePrice = avgPrice
}

func (s *DataService) RefreshPoolInfo() {
	poolBrBalance, err := data.GetTokenBalance(global.BscClient, constx.BrAddress, constx.BrPoolAddress)
	if err != nil {
		global.Logger.Error("获取流动性池异常")
		global.Logger.Error(err.Error())
		return
	}
	s.PoolInfo.BrBalance = poolBrBalance

	poolUsdtBalance, err := data.GetTokenBalance(global.BscClient, constx.UsdtAddress, constx.BrPoolAddress)
	if err != nil {
		global.Logger.Error("获取流动性池异常")
		global.Logger.Error(err.Error())
		return
	}
	s.PoolInfo.UsdtBalance = poolUsdtBalance

	s.SaveData(constx.RedisKeyPoolBR, fmt.Sprintf("%s", s.PoolInfo.BrBalance.String()), time.Now())
	s.SaveData(constx.RedisKeyPoolUsdt, fmt.Sprintf("%s", s.PoolInfo.UsdtBalance.String()), time.Now())

	avgPoolBr, err := s.GetData(constx.RedisKeyPoolBR)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}
	s.PoolInfo.AvgBrBalance = avgPoolBr

	avgPoolUsdt, err := s.GetData(constx.RedisKeyPoolUsdt)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}
	s.PoolInfo.AvgUsdtBalance = avgPoolUsdt

}

// 期现价差
func (s *DataService) GetBasisPct() decimal.Decimal {
	return s.BrFuturePrice.Sub(s.BrPrice).Div(s.BrPrice)
}

// 现货与均价的差
func (s *DataService) GetPriceToAvgSpreadPct() decimal.Decimal {
	return s.BrPrice.Sub(s.AvgBrPrice).Div(s.AvgBrPrice)
}

func (s *DataService) GetBrPoolBalanceLow() bool {
	rate := s.PoolInfo.BrBalance.Sub(s.PoolInfo.AvgBrBalance).Div(s.PoolInfo.AvgBrBalance)
	if rate.LessThanOrEqual(decimal.NewFromFloat(-0.7)) {
		return true
	}

	return false
}

func (s *DataService) GetUsdtPoolBalanceLow() bool {

	rate := s.PoolInfo.UsdtBalance.Sub(s.PoolInfo.AvgUsdtBalance).Div(s.PoolInfo.AvgUsdtBalance)

	if rate.LessThanOrEqual(decimal.NewFromFloat(-0.7)) {
		return true
	}

	return false
}
