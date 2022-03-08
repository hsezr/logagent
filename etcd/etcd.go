package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"logagent/common"
	"logagent/tailfile"
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
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
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
	fmt.Println(string(ret.Value))
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed, err:%v", err)
		return
	}

	return
}

//监控etcd中日志收集项配置变化的函数
func WatchConf(key string) {
	watchCh := cli.Watch(context.Background(), key)
	for wresp := range watchCh {
		logrus.Info("get new conf from etcd!")
		for _, evt := range wresp.Events {
			fmt.Printf("type:%s key:%s value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
			var newConf []common.CollectEntry
			if evt.Type == clientv3.EventTypeDelete {
				tailfile.SendNewConf(newConf)
				continue
			}
			err := json.Unmarshal(evt.Kv.Value, &newConf)
			if err != nil {
				logrus.Errorf("json unmarshal new conf failed, err:%v", err)
				continue
			}
			tailfile.SendNewConf(newConf)
		}
	}
}
