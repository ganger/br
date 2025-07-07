package util

import (
	"br-trade/bootstrap"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTokenPrice(t *testing.T) {
	bootstrap.InitBscClient()

	ts, buy1, sell1, err := GetTokenPrice("BRUSDT")
	assert.Nil(t, err)
	fmt.Println(ts, buy1, sell1)
}
