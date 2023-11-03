package config

type AppConf struct {
	KafkaConf `ini:"kafka"`
	//TailLogConf `ini:"tail_log"`
	EtcdConf `ini:"etcd"`
	ESConf   `ini:"es"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
	MaxSize int    `ini:"max_size"`
}

//type TailLogConf struct {
//	LogPath string `ini:"log_path"`
//}

type EtcdConf struct {
	Address       string `ini:"address"`
	Timeout       int    `ini:"timeout"`
	CollectLogKey string `ini:"collect_log_key"`
}
type ESConf struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chan_size"`
	Nums     int    `ini:"nums"`
}

func LoadConf(confType, filename string) (err error) {
	//kafka
	return
}
