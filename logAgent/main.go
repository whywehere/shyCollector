package main

import (
	"gopkg.in/ini.v1"
	"log/slog"
	"shyCollector/logAgent/config"
	"shyCollector/logAgent/etcd"
	"shyCollector/logAgent/kafka"
	"shyCollector/logAgent/tailLog"
	"time"
)

var (
	cfg = new(config.AppConf)
)

func main() {
	if err := ini.MapTo(&cfg, "C:\\Users\\19406\\Desktop\\go\\shyCollector\\logAgent\\config\\config.ini"); err != nil {
		panic(err)
	}

	// initialize kafka
	addr := []string{cfg.KafkaConf.Address}
	if err := kafka.Init(addr); err != nil {
		panic(err)
	}
	slog.Info("initialize kafka successfully")

	// initialize tailLog
	if err := tailLog.Init(cfg.LogPath); err != nil {
		panic(err)
	}
	slog.Info("initialize tailLog successfully")

	// initialize etcd
	etcdAddr := []string{cfg.EtcdConf.Address}
	if err := etcd.Init(etcdAddr, time.Duration(cfg.EtcdConf.Timeout)*time.Second); err != nil {
		panic(err)
	}
	slog.Info("initialize etcd successfully")
	go run()
	select {}
}

func run() {
	for {
		select {
		case line := <-tailLog.ReadChan():
			kafka.SendToKafka(cfg.Topic, line.Text)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}
