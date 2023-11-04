package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log/slog"
	"os"
	"os/signal"
	"shyCollector/logAgent/config"
	"shyCollector/logAgent/etcd"
	"shyCollector/logAgent/kafka"
	"shyCollector/logAgent/tailLog"
	"syscall"
	"time"
)

func main() {

	// 启动信号监听
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	cfg := new(config.AppConf)
	if err := ini.MapTo(&cfg, "C:\\Users\\19406\\Desktop\\go\\shyCollector\\config\\config.ini"); err != nil {
		panic(err)
	}

	// start kafka
	addr := []string{cfg.KafkaConf.Address}
	if err := kafka.Start(addr, cfg.KafkaConf.MaxSize); err != nil {
		panic(err)
	}
	slog.Info("start kafka successfully")

	// initialize etcd
	etcdAddr := []string{cfg.EtcdConf.Address}
	logEntryConf, err := etcd.Start(etcdAddr, time.Duration(cfg.EtcdConf.Timeout)*time.Second, cfg.CollectLogKey)
	if err != nil {
		panic(err)
	}
	slog.Info("initialize etcd successfully")

	// initialize tailLog
	if err := tailLog.Start(logEntryConf, cfg.CollectLogKey); err != nil {
		panic(err)
	}
	slog.Info("start tailLog successfully")

	// 阻塞等待信号
	<-sigChan
	fmt.Println("Received signal. Exiting gracefully.")
	os.Exit(0)
}
