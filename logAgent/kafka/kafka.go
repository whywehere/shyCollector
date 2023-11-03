package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"time"
)

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addrs []string, maxSize int) (err error) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出⼀个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		slog.Error("sarama.NewSyncProducer()", "Error", err)
		return
	}
	logDataChan = make(chan *logData, maxSize)
	go sendToKafka()
	return
}

func sendToKafka() {
	for {
		select {
		case msg := <-logDataChan:
			message := &sarama.ProducerMessage{
				Topic: msg.topic,
				Value: sarama.StringEncoder(msg.data),
			}

			// 发送消息
			pid, offset, err := client.SendMessage(message)
			if err != nil {
				slog.Error("kafka send message failed", "Error", err)
				return
			}
			slog.Info(fmt.Sprintf("pid:%v offset:%v\n", pid, offset))
		default:
			time.Sleep(time.Second)
		}
	}
	// 构造⼀个消息

}

// SendToChan 将外部消息存入内部channel
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}
