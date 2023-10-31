package etcd

import (
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
