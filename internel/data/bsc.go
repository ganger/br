package data

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

const balanceOfMethod = "0x70a08231"

func GetTokenBalance(client *ethclient.Client, tokenAddress, holderAddress string) (decimal.Decimal, error) {
	// 构造调用数据
	// 将holderAddress转换为正确的格式（去掉0x前缀，左补零到32字节）
	paddedAddress := common.LeftPadBytes(common.FromHex(holderAddress), 32)

	// 构造完整的调用数据：balanceOf方法签名 + 参数
	data := append(common.FromHex(balanceOfMethod), paddedAddress...)

	// 调用合约
	addr := common.HexToAddress(tokenAddress)
	callMsg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to call contract: %v", err)
	}

	// 解析结果
	balance := new(big.Int)
	balance.SetBytes(result)

	d := decimal.NewFromBigInt(balance, 0)
	divisor := decimal.New(1, 18) // 相当于 10^decimals
	d = d.Div(divisor)
	return d, nil
}
