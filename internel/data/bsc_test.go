package data

import (
	"br-trade/bootstrap"
	"br-trade/constx"
	"br-trade/global"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContractBalance(t *testing.T) {
	bootstrap.InitBscClient()

	result, err := GetTokenBalance(global.BscClient, constx.BrAddress, constx.BrPoolAddress)
	assert.Nil(t, err)
	fmt.Println(result)

	result2, err := GetTokenBalance(global.BscClient, constx.UsdtAddress, constx.BrPoolAddress)
	assert.Nil(t, err)
	fmt.Println(result2)
}
