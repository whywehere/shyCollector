package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestEtcdOptions(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("EtcdCli started...")
	defer cli.Close()

	//Put
	if _, err = cli.Put(context.Background(), "/logAgent/collect_config", "[{\"path\": \"C:\\\\Users\\\\19406\\\\Desktop\\\\go\\\\shyCollector\\\\logAgent\\\\logtest1.log\", \"topic\": \"web_log\"},\n{\"path\": \"C:\\\\Users\\\\19406\\\\Desktop\\\\go\\\\shyCollector\\\\logAgent\\\\logtest2.log\", \"topic\": \"redis_log\"}]"); err != nil {
		t.Fatal(err)
	}
	// Get
	resp, err := cli.Get(context.Background(), "/logAgent/collect_config")
	if err != nil {
		t.Fatal(err)
	}
	for _, kv := range resp.Kvs {
		t.Log(string(kv.Key), " ", string(kv.Value))
	}
}

func TestWatch(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ch := cli.Watch(context.Background(), "/logAgent/collect_config")
	for resp := range ch {
		for _, evt := range resp.Events {
			t.Log(evt)
		}
	}

}
