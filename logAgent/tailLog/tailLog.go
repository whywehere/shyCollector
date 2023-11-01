package tailLog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"log/slog"
	"shyCollector/logAgent/kafka"
)

var (
	tailObj *tail.Tail
	LogChan chan string
)

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	ctx      context.Context
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:  path,
		topic: topic,
		ctx:   context.Background(),
	}
	tailObj.Init()
	return
}

func (t *TailTask) Init() {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true}
	tailObj, err := tail.TailFile(t.path, config)
	if err != nil {
		panic(fmt.Sprintf("tailTask init failed: %v\n", err))
	}
	t.instance = tailObj

	go t.run()
}

func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines:
			kafka.SendToChan(t.topic, line.Text)
		case <-t.ctx.Done():
			slog.Info(fmt.Sprintf("%s_%s tailTask is done", t.topic, t.path))
			return

		}
	}
}
