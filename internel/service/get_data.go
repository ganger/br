package service

import (
	"br-trade/constx"
	"br-trade/global"
	"br-trade/internel/data"
	"fmt"
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
}
