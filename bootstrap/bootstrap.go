package bootstrap

import (
	"br-trade/conf"
	"br-trade/constx"
	"br-trade/global"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/yaml.v3"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/redis/go-redis/v9"
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
	fmt.Println(fmt.Sprintf("%+v", global.Config))
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

	global.BinanceFuturesClient = binance.NewFuturesClient(global.Config.Binance.ApiKey, global.Config.Binance.SecretKey)

}

func InitRedis() {

	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 密码（如果没有密码则留空）
		DB:       0,                // 默认数据库
	})

	// 测试连接
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	global.RedisClient = rdb
	fmt.Println("redis连接成功:", pong)
}
