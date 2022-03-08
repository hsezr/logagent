package tailfile

import (
	"context"
	"logagent/kafka"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	tt := &tailTask{
		path:  path,
		topic: topic,
	}
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tt.tObj, _ = tail.TailFile(path, cfg)
	tt.ctx, tt.cancel = context.WithCancel(context.Background())
	return tt
}

func (t *tailTask) run() {
	//读取日志发往kafka
	logrus.Infof("collect for path:%s is running...", t.path)
	for {
		select {
		case <-t.ctx.Done():
			logrus.Infof("path:%s is stopping...", t.path)
			return
		case line, ok := <-t.tObj.Lines:
			if !ok {
				logrus.Warn("tail file close reopen, filename: %s\n", t.path)
				time.Sleep(time.Second)
				continue
			}
			if len(line.Text) == 0 {
				continue
			}
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			logrus.Info("line.Text: %v", line.Text)
			kafka.MsgChan <- msg
		}
		
	}
}
