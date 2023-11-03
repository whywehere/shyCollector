package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log/slog"
	"strings"
	"time"
)

type LogData struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}

var (
	Client *elastic.Client
	ch     chan *LogData
)

func Start(addr string, nums, chanSize int) (err error) {
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}
	Client, err = elastic.NewClient(elastic.SetURL(addr))
	if err != nil {
		return
	}
	ch = make(chan *LogData, chanSize)
	for i := 0; i < nums; i++ {
		go sendToES()
	}

	return
}

func sendToES() {
	for {
		select {
		case msg := <-ch:
			do, err := Client.Index().Index(msg.Topic).BodyJson(msg).Do(context.Background())
			if err != nil {
				slog.Error("", "Error", err)
				continue
			}
			slog.Info(fmt.Sprintf("Index : %s, type : %s", do.Index, do.Type))
		default:
			time.Sleep(time.Second)

		}
	}

}

func SendToChan(msg *LogData) {
	ch <- msg
}
