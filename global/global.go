package global

import (
	"br-trade/conf"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config                conf.Config
	Logger                *zap.Logger
	DB                    *gorm.DB
	BscClient             *ethclient.Client
	BinanceFuturesClient  *futures.Client
	BinanceFuturesClient2 *futures.Client
	RedisClient           *redis.Client
)
