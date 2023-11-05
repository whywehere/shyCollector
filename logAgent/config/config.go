package config

type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
	MaxSize int    `ini:"max_size"`
}

type EtcdConf struct {
	Address       string `ini:"address"`
	Timeout       int    `ini:"timeout"`
	CollectLogKey string `ini:"collect_log_key"`
}
