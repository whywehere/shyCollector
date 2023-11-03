package main

import (
	"gopkg.in/ini.v1"
	"shyCollector/config"
	"shyCollector/logTransfer/es"
	"shyCollector/logTransfer/kafka"
)

var cfg = new(config.AppConf)

func main() {
	if err := ini.MapTo(cfg, "C:\\Users\\19406\\Desktop\\go\\shyCollector\\config\\config.ini"); err != nil {
		panic(err)
	}

	if err := es.Start(cfg.ESConf.Address, cfg.ESConf.Nums, cfg.ESConf.ChanSize); err != nil {
		panic(err)
	}

	if err := kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic); err != nil {
		panic(err)
	}

}
