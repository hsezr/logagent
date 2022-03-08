package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

func Init(address []string, chanSize int64) (err error) {
	//生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka: producer closed, err:", err)
		return
	}
	MsgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg()
	return
}

func sendMsg() {
	for {
		select {
		case msg := <-MsgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed, err:", err)
				return
			}
			logrus.Info("send msg to kafka success. pid:%v offset:%v", pid, offset)
		}
	}
}
