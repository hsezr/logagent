package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"logagent/common"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

//etcd相关操作
var (
	cli *clientv3.Client
)

func Init(addr []string) (err error) {
	cli, err = clientv3.New(
		clientv3.Config{
			Endpoints:   addr,
			DialTimeout: time.Second * 5,
		})
	if err != nil {
		fmt.Printf("")
		return
	}
	return
}

//拉去日志收集配置项的函数
func GetConf(key string) (collectEntryList []*common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	resp, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key:%s failed, err: %v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Errorf("get conf from etcd by key len = 0 by :%s failed, err: %v", key, err)
		return
	}
	ret := resp.Kvs[0]
	fmt.Println(ret.Value)
	err = json.Unmarshal(ret.Value, collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed, err:%v", err)
		return
	}

	return
}
