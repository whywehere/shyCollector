package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestEtcdOptions(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("EtcdCli started...")
	defer cli.Close()

}
