package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var etcdCli *clientv3.Client

func Init(addr []string, timeout time.Duration) (err error) {
	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})
	return
}

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func GetConf(key string) (value []*LogEntry, err error) {
	resp, err := etcdCli.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		//slog.Info(fmt.Sprintf("%s: %s", kv.Key, kv.Value))
		if err := json.Unmarshal(kv.Value, &value); err != nil {
			return nil, err
		}
	}
	return
}
