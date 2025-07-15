package service

import (
	"br-trade/constx"
	"br-trade/global"
	"br-trade/internel/data"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type DataService struct {
	BrPrice       decimal.Decimal
	BrFuturePrice decimal.Decimal
	PoolInfo      PoolInfo

	ShutDown bool
}

type PoolInfo struct {
	BrBalance   decimal.Decimal
	UsdtBalance decimal.Decimal
}

func NewDataService() *DataService {
	return &DataService{}
}

func (s *DataService) Run() {
	s.Init()

	go func() {
		for {
			if s.ShutDown {
				return
			}
			s.RefreshBrPrice()
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			if s.ShutDown {
				return
			}
			s.RefreshBrFuturePrice()
			time.Sleep(1 * time.Second)
		}

	}()

	go func() {
		for {
			if s.ShutDown {
				return
			}
			s.RefreshPoolInfo()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			if s.ShutDown {
				return
			}
			s.PushWx()
			time.Sleep(1 * time.Hour)
		}
	}()

	for {

		if s.ShutDown {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *DataService) Init() {
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
	global.Logger.Info("服务启动成功")
}

func (s *DataService) PushWx() {

	msg := fmt.Sprintf("BR期货价格:%s\n", s.BrFuturePrice.Round(6).String())
	averageBrFuturePrice, err := s.GetData(constx.RedisKeyBrFuturePrice)
	if err == nil {
		msg = msg + fmt.Sprintf("BR期货平均价格:%s\n", averageBrFuturePrice.Round(6).String())
	}

	msg = msg + fmt.Sprintf("BR现货价格:%s\n", s.BrPrice.Round(6).String())
	averageBrPrice, err := s.GetData(constx.RedisKeyBrPrice)
	if err == nil {
		msg = msg + fmt.Sprintf("BR现货平均价格:%s\n", averageBrPrice.Round(6).String())
	}

	msg = msg + fmt.Sprintf("流动性池子BR余额:%sM\n",
		s.PoolInfo.BrBalance.Div(decimal.NewFromInt(1000000)).Round(2).String())
	averagePoolBr, err := s.GetData(constx.RedisKeyPoolBR)
	if err == nil {
		msg = msg + fmt.Sprintf("BR池子平均数量:%sM\n",
			averagePoolBr.Div(decimal.NewFromInt(1000000)).Round(2).String())
	}

	msg = msg + fmt.Sprintf("流动性池子USDT余额:%sM\n",
		s.PoolInfo.UsdtBalance.Div(decimal.NewFromInt(1000000)).Round(2).String())
	averagePoolUsdt, err := s.GetData(constx.RedisKeyPoolUsdt)
	if err == nil {
		msg = msg + fmt.Sprintf("USDT池子平均数量:%sM\n",
			averagePoolUsdt.Div(decimal.NewFromInt(1000000)).Round(2).String())
	}

	poolLiquidity := s.PoolInfo.BrBalance.Mul(s.BrPrice).Add(s.PoolInfo.UsdtBalance)
	msg = msg + fmt.Sprintf("流动性总金额:%sM\n",
		poolLiquidity.Div(decimal.NewFromInt(1000000)).Round(2).String())
	global.Logger.Info(msg)
	//util.PushWX(global.Config.Wx.MessagePushUrl, msg)
}

func (s *DataService) Stop() {
	s.ShutDown = true
}

/*
1.现货没变，pool流动性正常，期货价格突变
2.现货价格波动，流动性减少（非必须）
*/
func (s *DataService) CheckPosition() {

	if s.BrPrice.Sub(s.BrFuturePrice).Div(s.BrPrice).GreaterThan(decimal.NewFromFloat(0.01)) {

	}
}
