package tailLog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"log/slog"
	"shyCollector/logAgent/kafka"
)

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	tailCtx  context.Context
}

func NewTailTask(path, topic string) (*TailTask, error) {

	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true}
	tailX, err := tail.TailFile(path, config)
	if err != nil {
		return nil, err
	}
	t := &TailTask{
		path:    path,
		topic:   topic,
		tailCtx: context.Background(),
	}
	t.instance = tailX
	go t.serve()
	return t, nil
}

func (t *TailTask) serve() {
	for {
		select {
		case line := <-t.instance.Lines:
			kafka.SendToChan(t.topic, line.Text)
		case <-t.tailCtx.Done():
			slog.Info(fmt.Sprintf("%s_%s tailTask is done", t.topic, t.path))
			return
		}
	}
}
