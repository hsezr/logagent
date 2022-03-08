package tailfile

import (
	"logagent/common"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

type tailTask struct {
	path string
	topic string
	tObj *tail.Tail
}

func Init(allConf []*common.CollectEntry) (err error) {
	// allConf
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	TailObj, err = tail.TailFile(filename, config)
	if err != nil {
		logrus.Error("tailfile path: %s failed, err:%v\n", filename, err)
		return
	}
	return
}
