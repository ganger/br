package service

import (
	"br-trade/constx"
	"br-trade/global"
	"br-trade/internel/data"
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

type DataService struct {
	BrPrice          decimal.Decimal
	AvgBrPrice       decimal.Decimal
	BrFuturePrice    decimal.Decimal
	AvgBrFuturePrice decimal.Decimal
	PoolInfo         PoolInfo

	ShutDown bool
}

type PoolInfo struct {
	BrBalance      decimal.Decimal
	AvgBrBalance   decimal.Decimal
	UsdtBalance    decimal.Decimal
	AvgUsdtBalance decimal.Decimal
}

func NewDataService() *DataService {
	return &DataService{}
}

func (s *DataService) Run() {
	s.Init()

	go func() {
		for {
			if s.ShutDown {
				break
			}
			s.RefreshBrPrice()
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			if s.ShutDown {
				break
			}
			s.RefreshBrFuturePrice()
			time.Sleep(1 * time.Second)
		}

	}()

	go func() {
		for {
			if s.ShutDown {
				break
			}
			s.RefreshPoolInfo()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			if s.ShutDown {
				break
			}
			s.CheckPosition()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if s.ShutDown {
				break
			}
			s.PushWx()
			time.Sleep(1 * time.Hour)
		}
	}()

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

	basisPct := s.GetBasisPct()
	msg = msg + fmt.Sprintf("期现差价:%s%%\n", basisPct.Mul(decimal.NewFromInt(100)).Round(4).String())

	pct := s.GetPriceToAvgSpreadPct()
	if err == nil {
		msg = msg + fmt.Sprintf("现货与均价偏差:%s%%\n", pct.Mul(decimal.NewFromInt(100)).Round(4).String())
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

	priceToAvgSpreadPct := s.GetPriceToAvgSpreadPct()

	//现货涨幅超过1%
	if priceToAvgSpreadPct.GreaterThanOrEqual(decimal.NewFromFloat(0.01)) {
		//池子br数量降低
		isPoolLow := s.GetBrPoolBalanceLow()
		if isPoolLow {
			//现货准备上涨，开仓
			s.CreateOrder(futures.SideTypeBuy)
		}
	}

	//现货跌幅超过1%
	if priceToAvgSpreadPct.LessThanOrEqual(decimal.NewFromFloat(-0.01)) {
		//池子usdt数量降低
		isPoolLow := s.GetUsdtPoolBalanceLow()
		if isPoolLow {
			//现货准备下跌，开仓
			s.CreateOrder(futures.SideTypeSell)
		}
	}
}

func (s *DataService) CreateOrder(dir futures.SideType) {

	if !global.Config.App.IsPrd {
		global.Logger.Info("非prd，不下单")
		return
	}
	price0 := s.AvgBrPrice.Mul(decimal.NewFromFloat(0.96))
	price1 := s.AvgBrPrice.Mul(decimal.NewFromFloat(0.98))
	price2 := s.AvgBrPrice
	price3 := s.AvgBrPrice.Mul(decimal.NewFromFloat(1.02))
	price4 := s.AvgBrPrice.Mul(decimal.NewFromFloat(1.04))

	priceList := []decimal.Decimal{price0, price1, price2, price3, price4}
	if dir == futures.SideTypeSell {
		priceList = []decimal.Decimal{price4, price3, price2, price1, price0}
	}

	for _, price := range priceList {

		quantity := decimal.NewFromInt(4999).Div(price).Round(0)
		_, err := global.BinanceFuturesClient.NewCreateOrderService().
			Symbol(constx.BrFutureSymbol).
			Side(dir).
			Type(futures.OrderTypeLimit).
			TimeInForce(futures.TimeInForceTypeGTC).
			Price(price.Round(5).String()).
			Quantity(quantity.String()).
			Do(context.Background())
		if err != nil {
			global.Logger.Error(err.Error())
			continue
		}
		global.Logger.Info("下单成功",
			zap.String("价格", price.Round(5).String()),
			zap.String("数量", quantity.String()),
			zap.String("总价", price.Round(5).Mul(quantity).String()),
		)
	}

}
