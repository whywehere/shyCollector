package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log/slog"
	"shyCollector/config"
	"shyCollector/logAgent/etcd"
	"shyCollector/logAgent/kafka"
	"shyCollector/logAgent/tailLog"
	"shyCollector/utils"
	"sync"
	"time"
)

var (
	cfg = new(config.AppConf)
)

func main() {
	if err := ini.MapTo(&cfg, "C:\\Users\\19406\\Desktop\\go\\shyCollector\\config\\config.ini"); err != nil {
		panic(err)
	}

	// initialize kafka
	addr := []string{cfg.KafkaConf.Address}
	if err := kafka.Init(addr, cfg.KafkaConf.MaxSize); err != nil {
		panic(err)
	}
	slog.Info("initialize kafka successfully")

	// initialize etcd
	etcdAddr := []string{cfg.EtcdConf.Address}
	if err := etcd.Init(etcdAddr, time.Duration(cfg.EtcdConf.Timeout)*time.Second); err != nil {
		panic(err)
	}
	slog.Info("initialize etcd successfully")
	ip, err := utils.GetOutBoundIp()
	if err != nil {
		panic(err)
	}
	collectLogKEY := fmt.Sprintf(cfg.CollectLogKey, ip)
	logEntryConf, err := etcd.GetConf(collectLogKEY)
	if err != nil {
		slog.Error("etcd.GetConf", "Error", err)
		return
	}

	slog.Info(fmt.Sprintf("get etcdConf successfully, %v\n", logEntryConf))
	// initialize tailLog
	tailLog.Init(logEntryConf)
	slog.Info("initialize tailLog successfully")

	newConfChan := tailLog.NewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(cfg.CollectLogKey, newConfChan)
	wg.Wait()
}
