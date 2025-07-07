package service

import (
	"br-trade/constx"
	"br-trade/global"
	"br-trade/internel/util"
	"fmt"
	"github.com/shopspring/decimal"
)

func (s DataService) Run() {
	s.Init()
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
	_, buy, _, err := util.GetFuturePrice(constx.BrAddress)
	if err != nil {
		panic(err)
	}
	s.BrFuturePrice = buy

	price, err := util.GetTokenPriceFromBinance(constx.BrAddress)
	if err != nil {
		panic(err)
	}
	s.BrPrice = price

	poolBrBalance, err := util.GetTokenBalance(global.BscClient, constx.BrAddress, constx.BrPoolAddress)
	if err != nil {
		panic(err)
	}
	s.PoolInfo.BrBalance = poolBrBalance

	poolUsdtBalance, err := util.GetTokenBalance(global.BscClient, constx.UsdtAddress, constx.BrPoolAddress)
	if err != nil {
		panic(err)
	}
	s.PoolInfo.UsdtBalance = poolUsdtBalance
	s.PushWx()
}

func (s DataService) PushWx() {

	msg := fmt.Sprintf("服务启动成功\n===========================\n")
	msg = msg + fmt.Sprintf("BR期货价格:%s\n", s.BrFuturePrice.Round(6).String())
	msg = msg + fmt.Sprintf("BR现货价格:%s\n", s.BrPrice.Round(6).String())

	msg = msg + fmt.Sprintf("流动性池子BR余额:%s\n", s.PoolInfo.BrBalance.Round(6).String())
	msg = msg + fmt.Sprintf("流动性池子USDT余额:%s\n", s.PoolInfo.UsdtBalance.Round(6).String())
	poolLiquidity := s.PoolInfo.BrBalance.Mul(s.BrPrice).Add(s.PoolInfo.UsdtBalance)
	msg = msg + fmt.Sprintf("流动性总金额:%s\n", poolLiquidity.Round(6).String())
	util.PushWX(global.Config.Wx.MessagePushUrl, msg)
}

func (s DataService) Stop() {
}
