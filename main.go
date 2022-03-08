package main

import (
	"fmt"
	"logagent/common"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	Chansize int64  `ini:"chan_size"`
}


type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func main() {
	//获取本机IP
	ip, err := common.GetOutboundIP()
	if err != nil {
		logrus.Errorf("get ip failed, err: %v", err)
		return
	}
	var configObj = new(Config)
	// cfg, err := ini.Load("./conf/config.ini")
	// if err != nil {
	// 	logrus.Error("load config failed, err: %v", err)
	// 	return
	// }
	// kafkaAddr := cfg.Section("kafka").Key("address").String()
	// fmt.Println(kafkaAddr)
	err = ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config failed, err:%v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)

	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.Chansize)
	if err != nil {
		logrus.Errorf("init kafka failed, err: %v", err)
		return
	}
	logrus.Info("init kafka success")

	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed, err:%v", err)
		return
	}
	//从etcd中拉取要收集的日志的配置项
	configObj.CollectKey = fmt.Sprintf(configObj.CollectKey, ip)
	allConf, err := etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("getconf from etcd failed, err:%v", err)
		return
	}
	fmt.Println(allConf)

	go etcd.WatchConf(configObj.EtcdConfig.CollectKey)

	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tail failed, err: %v", err)
		return
	}
	logrus.Info("init tailfile success!")

	if err != nil {
		logrus.Error()
	}

	select {

	}
}