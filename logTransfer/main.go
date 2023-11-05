package main

import (
	"flag"
	"gopkg.in/ini.v1"
	"shyCollector/logTransfer/config"
	"shyCollector/logTransfer/es"
	"shyCollector/logTransfer/kafka"
)

var cfg = new(config.AppConf)

func main() {
	if err := ini.MapTo(cfg, "C:\\Users\\19406\\Desktop\\go\\shyCollector\\config\\config.ini"); err != nil {
		panic(err)
	}
	flag.StringVar(&cfg.Topic, "topic", cfg.Topic, "")
	flag.StringVar(&cfg.KafkaConf.Address, "kafka_addr", cfg.KafkaConf.Address, "")
	flag.StringVar(&cfg.ESConf.Address, "es_addr", cfg.EtcdConf.Address, "")
	flag.Parse()
	if err := es.Start(cfg.ESConf.Address, cfg.ESConf.Nums, cfg.ESConf.ChanSize); err != nil {
		panic(err)
	}

	if err := kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic); err != nil {
		panic(err)
	}

}
