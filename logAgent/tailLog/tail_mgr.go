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
		tailTask := NewTailTask(entry.Path, entry.Topic)
		key := fmt.Sprintf("%s_%s", entry.Topic, entry.Path)
		tskMgr.tskMap[key] = tailTask
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
					tailObj := NewTailTask(entry.Path, entry.Topic)
					t.tskMap[confKey] = tailObj
				}
			}
			// stop the deleted tasks
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
					t.tskMap[key].ctx.Done()
				}
			}

		default:
			time.Sleep(time.Second)
		}
	}
}

func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
