package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
)

type GateTokenPriceResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		CurrentPrice float64 `json:"current_price"`
		Time         int64   `json:"time"`
		Symbol       string  `json:"symbol"`
		Decimal      int     `json:"decimal"`
	} `json:"data"`
}

func GetTokenPriceFromGate(address string) (price decimal.Decimal, err error) {
	url := fmt.Sprintf("https://apipro.gateweb3.cc/web3api/v2/token/current_price?chain=bsc&address=%s&native_coin=false", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gtweb3-Device-Type", "3")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get token price from gate web3 error: %s", resp.Status)
		return
	}

	var result GateTokenPriceResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	price = decimal.NewFromFloat(result.Data.CurrentPrice)
	return
}

type BinanceTokenPriceResponse struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"` // 使用interface{}因为可能是null
	MessageDetail interface{} `json:"messageDetail"`
	Data          KlineData   `json:"data"`
	Success       bool        `json:"success"`
}
type KlineData struct {
	KlineInfos [][]string `json:"klineInfos"` // 每个K线数据是一个字符串数组
	Decimals   int        `json:"decimals"`
}

func GetTokenPriceFromBinance(address string) (price decimal.Decimal, err error) {
	url := fmt.Sprintf("https://www.binance.com/bapi/defi/v1/public/alpha-trade/agg-klines?chainId=56&interval=1s&limit=1&tokenAddress=%s", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gtweb3-Device-Type", "3")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get token price from gate web3 error: %s", resp.Status)
		return
	}

	var result BinanceTokenPriceResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	if !result.Success || result.Code != "000000" {
		err = errors.New("BinanceTokenPriceResponse fail")
		return
	}
	if len(result.Data.KlineInfos) < 1 || len(result.Data.KlineInfos[0]) < 7 {
		err = fmt.Errorf("BinanceTokenPriceResponse error: %v", result.Data)
		return
	}
	price, err = decimal.NewFromString(result.Data.KlineInfos[0][4])
	if err != nil {
		return
	}
	return
}
