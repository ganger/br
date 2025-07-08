package service

import (
	"br-trade/constx"
	"br-trade/global"
	"br-trade/internel/data"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

func (s DataService) Run() {
	s.Init()

	go func() {
		s.RefreshBrPrice()
		time.Sleep(1 * time.Second)
	}()

	go func() {
		s.RefreshBrFuturePrice()
		time.Sleep(1 * time.Second)
	}()

	go func() {
		s.RefreshPoolInfo()
		time.Sleep(1 * time.Second)
	}()

	go func() {
		s.PushWx()
		time.Sleep(10 * time.Second)
	}()
}

type DataService struct {
	BrPrice       decimal.Decimal
	BrFuturePrice decimal.Decimal
	PoolInfo      PoolInfo
}

type PoolInfo struct {
	BrBalance   decimal.Decimal
	UsdtBalance decimal.Decimal
}

func NewDataService() *DataService {
	return &DataService{}
}

func (s DataService) Init() {
	_, buy, _, err := data.GetFuturePrice(constx.BrFutureSymbol)
	if err != nil {
		panic(err)
	}
	s.BrFuturePrice = buy

	price, err := data.GetTokenPriceFromBinance(constx.BrAddress)
	if err != nil {
		panic(err)
	}
	s.BrPrice = price

	poolBrBalance, err := data.GetTokenBalance(global.BscClient, constx.BrAddress, constx.BrPoolAddress)
	if err != nil {
		panic(err)
	}
	s.PoolInfo.BrBalance = poolBrBalance

	poolUsdtBalance, err := data.GetTokenBalance(global.BscClient, constx.UsdtAddress, constx.BrPoolAddress)
	if err != nil {
		panic(err)
	}
	s.PoolInfo.UsdtBalance = poolUsdtBalance
	s.PushWx()
}

func (s DataService) PushWx() {

	msg := fmt.Sprintf("服务启动成功\n=======================\n")
	msg = msg + fmt.Sprintf("BR期货价格:%s\n", s.BrFuturePrice.Round(6).String())
	msg = msg + fmt.Sprintf("BR现货价格:%s\n", s.BrPrice.Round(6).String())

	msg = msg + fmt.Sprintf("流动性池子BR余额:%s\n", s.PoolInfo.BrBalance.Round(6).String())
	msg = msg + fmt.Sprintf("流动性池子USDT余额:%s\n", s.PoolInfo.UsdtBalance.Round(6).String())
	poolLiquidity := s.PoolInfo.BrBalance.Mul(s.BrPrice).Add(s.PoolInfo.UsdtBalance)
	msg = msg + fmt.Sprintf("流动性总金额:%s\n", poolLiquidity.Round(6).String())
	fmt.Println(msg)
	//util.PushWX(global.Config.Wx.MessagePushUrl, msg)
}

func (s DataService) RefreshBrPrice() {
	price, err := data.GetTokenPriceFromBinance(constx.BrAddress)
	if err != nil {
		global.Logger.Error("获取现货价格异常")
		global.Logger.Error(err.Error())
		return
	}
	s.BrPrice = price
}

func (s DataService) RefreshBrFuturePrice() {
	_, buy, _, err := data.GetFuturePrice(constx.BrFutureSymbol)
	if err != nil {
		global.Logger.Error("获取期货价格异常")
		global.Logger.Error(err.Error())
		return
	}
	s.BrFuturePrice = buy
}

func (s DataService) RefreshPoolInfo() {
	poolBrBalance, err := data.GetTokenBalance(global.BscClient, constx.BrAddress, constx.BrPoolAddress)
	if err != nil {
		if err != nil {
			global.Logger.Error("获取流动性池异常")
			global.Logger.Error(err.Error())
			return
		}
	}
	s.PoolInfo.BrBalance = poolBrBalance

	poolUsdtBalance, err := data.GetTokenBalance(global.BscClient, constx.UsdtAddress, constx.BrPoolAddress)
	if err != nil {
		if err != nil {
			global.Logger.Error("获取流动性池异常")
			global.Logger.Error(err.Error())
			return
		}
	}
	s.PoolInfo.UsdtBalance = poolUsdtBalance
}

func (s DataService) Stop() {
}
