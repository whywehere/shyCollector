package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"testing"
)

func TestKafkaProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出⼀个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	//addr := "127.0.0.1:9092"
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"       :9092"}, config)
	if err != nil {
		slog.Error("sarama.NewSyncProducer()", "Error", err)
		return
	}
	msg := &sarama.ProducerMessage{
		Value: sarama.StringEncoder("hello kafka test"),
		Topic: "web_log",
	}
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		return
	}
	slog.Info(fmt.Sprintf("pid:%v offset:%v\n", pid, offset))
}
