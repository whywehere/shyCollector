package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log/slog"
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

func WatchConf(key string, newConfChan chan<- []*LogEntry) {
	ch := etcdCli.Watch(context.Background(), key)
	for resp := range ch {
		var newConf []*LogEntry
		for _, evt := range resp.Events {
			if evt.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					slog.Error("Error unmarshalling ", "Error", err)
					continue
				}
			}
			slog.Info(fmt.Sprintf("Conf Changed: %v\n", newConf))
			newConfChan <- newConf
		}
	}
}
