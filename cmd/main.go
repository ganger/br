package main

import (
	"br-trade/bootstrap"
	"br-trade/internel/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bootstrap.InitConfig()
	bootstrap.InitLogger()
	bootstrap.InitBscClient()
	bootstrap.InitRedis()
	//	bootstrap.InitDB()
	svr := service.NewDataService()
	svr.Run()

	//等待中断信号以优雅地关闭服务器（设置 10 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server ...")
	time.Sleep(5 * time.Second)
	svr.Stop()
	log.Println("shutdown server done")
}
