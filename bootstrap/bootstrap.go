package bootstrap

import (
	"br-trade/conf"
	"br-trade/constx"
	"br-trade/global"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/yaml.v3"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitConfig() {
	data, err := os.ReadFile("./app.yml")
	if err != nil {
		panic(err)
	}

	// 解析 YAML 文件
	var config conf.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	global.Config = config
	if config.App.WalletAddress == "" {
		panic("wallet address list is empty")
	}
	if config.App.MessagePushUrl == "" {
		panic("message push url is empty")
	}
	if config.App.ReadApiKey == "" {
		panic("ReadApiKey is empty")
	}
	if config.App.ReadSecretKey == "" {
		panic("ReadSecretKey is empty")
	}
	initTime, err := time.ParseInLocation(time.DateOnly, config.App.InitTime, time.Local)
	if err != nil {
		panic(err)
	}
	fmt.Println("initTime:", initTime.Format(time.DateOnly))
	if config.App.ShortBTC > 0 {
		panic("ShortBTC > 0")
	}
	if config.App.IsPrd && (config.App.TradeApiKey == "" || config.App.TradeSecretKey == "") {
		panic("IsPrd && TradeApiKey or TradeSecretKey is empty")
	}
	fmt.Println("JlpBalance:", config.App.JlpBalance, ",ShortBTC:", config.App.ShortBTC)
}

func InitLogger() {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleOutput := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleOutput, zap.NewAtomicLevelAt(zap.InfoLevel))

	fileCfg := zap.NewProductionEncoderConfig()
	fileCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	fileEncoder := zapcore.NewJSONEncoder(fileCfg)
	logFile := "./logs/log"
	fileOutput, _ := rotatelogs.New(
		logFile+".%Y%m%d",
		rotatelogs.WithLinkName(logFile),
		rotatelogs.WithMaxAge(time.Hour*24*2),
	)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(fileOutput), zap.NewAtomicLevelAt(zap.DebugLevel))

	global.Logger = zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller())
	fmt.Println("init logger ok")
}

func InitBscClient() {
	client, err := ethclient.Dial(constx.BscEndPoint1)
	if err != nil {
		panic(err)
	}

	global.BscClient = client

}
