package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log/slog"
	"shyCollector/utils"
	"time"
)

var etcdCli *clientv3.Client

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Start(addr []string, timeout time.Duration, collectKey string) (value []*LogEntry, err error) {
	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})
	collectKey, err = utils.GetLocalCollectorKey(collectKey)
	if err != nil {
		return nil, err
	}
	resp, err := etcdCli.Get(context.Background(), collectKey)

	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		if err := json.Unmarshal(kv.Value, &value); err != nil {
			return nil, err
		}
	}
	return
}

func WatchConf(key string, newConfChan chan<- []*LogEntry) {
	watchChan := etcdCli.Watch(context.Background(), key)
	for resp := range watchChan {
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
