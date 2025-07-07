package util

import (
	"br-trade/constx"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTokenPriceFromGate(t *testing.T) {
	price, err := GetTokenPriceFromGate(constx.BrAddress)
	assert.Nil(t, err)
	fmt.Println(price)
}

func TestGetTokenPriceFromBinance(t *testing.T) {
	price, err := GetTokenPriceFromBinance(constx.BrAddress)
	assert.Nil(t, err)
	fmt.Println(price)
}
