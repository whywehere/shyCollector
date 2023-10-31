package config

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	TailLogConf `ini:"tail_log"`
	EtcdConf    `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type TailLogConf struct {
	LogPath string `ini:"log_path"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
}

func LoadConf(confType, filename string) (err error) {
	//kafka
	return
}
