package registry

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"micro/util/log"
	"testing"
	"time"
)

func TestEtcdRegistry(t *testing.T) {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.Get(context.Background(), "/micro/registry", clientv3.WithPrefix(), clientv3.WithKeysOnly())

	for _, v := range resp.Kvs {
		key := string(v.Key)
		log.Infof("Key:%s\n", key)
		resp2, _ := cli.Get(context.Background(), key)

		for _, v2 := range resp2.Kvs {
			log.Info(string(v2.Value))
		}
	}
}
