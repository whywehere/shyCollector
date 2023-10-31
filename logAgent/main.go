package main

import (
	"fmt"
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

	// initialize etcd
	etcdAddr := []string{cfg.EtcdConf.Address}
	if err := etcd.Init(etcdAddr, time.Duration(cfg.EtcdConf.Timeout)*time.Second); err != nil {
		panic(err)
	}
	slog.Info("initialize etcd successfully")

	logEntryConf, err := etcd.GetConf("/xxx")
	if err != nil {
		slog.Error("etcd.GetConf", "Error", err)
		return
	}
	slog.Info(fmt.Sprintf("get etcdConf successfully, %v\n", logEntryConf))
	tailLog.Init(logEntryConf)

	for _, entry := range logEntryConf {
		// initialize tailLog
		tailTask := tailLog.NewTailTask(entry.Path, entry.Topic)
		for {
			select {
			case line := <-tailObj.Lines:
				kafka.SendToKafka(entry.Topic, line.Text)
			}
		}
	}
	slog.Info("initialize tailLog successfully")
	select {}
}
