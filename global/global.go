package global

import (
	"br-trade/conf"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config    conf.Config
	Logger    *zap.Logger
	DB        *gorm.DB
	BscClient *ethclient.Client
)
