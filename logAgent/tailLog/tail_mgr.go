package tailLog

import (
	"fmt"
	"log/slog"
	"shyCollector/logAgent/etcd"
	"time"
)

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry:    logConf,
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry),
	}

	for _, entry := range tskMgr.logEntry {
		key := fmt.Sprintf("%s_%s", entry.Topic, entry.Path)
		task, err := NewTailTask(entry.Path, entry.Topic)
		if err != nil {
			slog.Error("Failed to create tail task")
			continue
		}
		tskMgr.tskMap[key] = task
	}

	go tskMgr.run()
}

func (t *tailLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			slog.Info(fmt.Sprintf("Conf Changed: %v\n", newConf))
			for _, entry := range newConf {
				confKey := fmt.Sprintf("%s_%s", entry.Topic, entry.Path)
				if _, ok := tskMgr.tskMap[confKey]; !ok {
					tailTask, err := NewTailTask(entry.Path, entry.Topic)
					if err != nil {
						slog.Error("Failed to create tail task")
						continue
					}
					t.tskMap[confKey] = tailTask
				}
			}
			// 停止并删除不再存在的任务
			for _, oldEntry := range t.logEntry {
				isDelete := true
				for _, newEntry := range newConf {
					if oldEntry.Topic == newEntry.Topic && oldEntry.Path == newEntry.Path {
						isDelete = false
						continue
					}
				}
				if isDelete {
					key := fmt.Sprintf("%s_%s", oldEntry.Topic, oldEntry.Path)
					t.tskMap[key].tailCtx.Done()
				}
			}
			t.logEntry = newConf
		default:
			time.Sleep(time.Second)
		}
	}
}

func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
