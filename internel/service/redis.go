package service

import (
	"br-trade/global"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func (s *DataService) SaveData(key string, value string, t time.Time) {

	keyTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), (t.Minute()/10)*10, 0, 0, time.Local)

	score := keyTime.Unix()
	scoreStr := strconv.Itoa(int(score))
	ctx := context.Background()

	count, err := global.RedisClient.ZCount(ctx, key, scoreStr, scoreStr).Result()
	if err != nil {
		global.Logger.Error("redis save data error", zap.Error(err))
	}
	if count > 0 {
		return
	}

	// 写入有序集合
	err = global.RedisClient.ZAdd(ctx, key, redis.Z{Score: float64(score), Member: value}).Err()
	if err != nil {
		return
	}

	// 清理24小时前的数据（确保只保留最近24h）
	minTime := float64(time.Now().Add(-24 * time.Hour).Unix())
	global.RedisClient.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%f", minTime))

}
