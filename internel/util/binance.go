package util

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
)

type OrderBookResponse struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	E            int64      `json:"E"`    // 消息时间（事件时间）
	T            int64      `json:"T"`    // 交易对时间
	Bids         [][]string `json:"bids"` // 买方深度 [价格, 数量]
	Asks         [][]string `json:"asks"` // 卖方深度 [价格, 数量]
}

func GetTokenPrice(token string) (price decimal.Decimal, err error) {
	url := fmt.Sprintf("https://www.binance.com/fapi/v1/depth?symbol=%s&limit=5", token)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	rsp := OrderBookResponse{}
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return
	}
	fmt.Println(rsp.E)
	fmt.Println(rsp.Bids[0][0])
	fmt.Println(rsp.Bids[1][0])
	return
}
