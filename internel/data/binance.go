package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"time"
)

type OrderBookResponse struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	E            int64      `json:"E"`    // 消息时间（事件时间）
	T            int64      `json:"T"`    // 交易对时间
	Bids         [][]string `json:"bids"` // 买方深度 [价格, 数量]
	Asks         [][]string `json:"asks"` // 卖方深度 [价格, 数量]
}

func GetFuturePrice(token string) (timeTs int64, buy1, sell1 decimal.Decimal, err error) {
	url := fmt.Sprintf("https://www.binance.com/fapi/v1/depth?symbol=%s&limit=5", token)
	resp, err := http.Get(url)
	if err != nil {
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
	timeTs = rsp.E
	if len(rsp.Bids) < 1 || len(rsp.Bids[0]) < 1 {
		err = errors.New("len(rsp.Bids) < 1 || len(rsp.Bids[0]) < 1")
		return
	}
	if len(rsp.Asks) < 1 || len(rsp.Asks[0]) < 1 {
		err = errors.New("len(rsp.Asks) < 1 || len(rsp.Asks[0]) < 1")
		return
	}
	now := time.Now()
	ts := time.UnixMilli(timeTs)
	if ts.Before(now.Add(-1 * time.Minute)) {
		err = fmt.Errorf("期货行情落后1分钟,价格时间:%s,当前时间:%s", ts.Format(time.DateTime), now.Format(time.DateTime))
		return
	}
	buy1, err = decimal.NewFromString(rsp.Bids[0][0])
	if err != nil {
		return
	}
	sell1, err = decimal.NewFromString(rsp.Asks[0][0])
	if err != nil {
		return
	}
	return
}
