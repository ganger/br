package util

import (
	"br-trade/bootstrap"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTokenPrice(t *testing.T) {
	bootstrap.InitBscClient()

	price, err := GetTokenPrice("BRUSDT")
	assert.Nil(t, err)

	fmt.Println(price)
}
