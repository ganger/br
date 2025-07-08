package data

import (
	"br-trade/bootstrap"
	"br-trade/constx"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTokenPrice(t *testing.T) {
	bootstrap.InitBscClient()

	ts, buy1, sell1, err := GetFuturePrice(constx.BrFutureSymbol)
	assert.Nil(t, err)
	fmt.Println(ts, buy1, sell1)
}
