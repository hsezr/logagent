package tailfile

import (
	"logagent/common"

	"github.com/sirupsen/logrus"
)

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       //所有的tailTask任务
	collectEntryList []common.CollectEntry      //所有配置项
	confChan         chan []common.CollectEntry //等待新配置的通道
}

var (
	ttMgr *tailTaskMgr
)

func Init(allConf []common.CollectEntry) (err error) {
	// allConf里存了若干个日志的 收集项
	//针对每一个日志收集项创建一个对应的tailObj
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic) //创建了一个日志收集的任务
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt //把创建的这个tailTask任务登记在册，方便管理
		go tt.run()
	}
	go ttMgr.watch() //在后台等新的配置来
	return
}

func (t *tailTaskMgr) watch() {
	for {
		newConf := <-t.confChan
		logrus.Infof("get new conf from etcd, NNNew conf: %v", newConf)
		//1.原来已经存在的不用动
		//2.原来没有的我要新创建一个tailtask任务
		//3.与拿来有的现在没有的要把tailTask停掉
		for _, conf := range newConf{
			if t.isExist(conf) {
				continue
			}

			tt := newTailTask(conf.Path, conf.Topic)
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			t.tailTaskMap[tt.path] = tt
			go tt.run()
		}
		//找出tailTaskMap中存在，但是newConf不存在的哪些tailTask，把他们停掉
		for key, task := range t.tailTaskMap {
			var found bool = false
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				logrus.Infof("the task collect path: %s need to stop.", task.path)
				delete(t.tailTaskMap, key)
				task.cancel()
			}
		}
	}
}

func (t *tailTaskMgr)isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
