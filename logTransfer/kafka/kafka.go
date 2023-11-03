package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"shyCollector/logTransfer/es"
)

func Init(addr []string, topic string) error {
	consumer, err := sarama.NewConsumer(addr, nil)
	if err != nil {
		return err
	}
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			return err
		}
		defer pc.AsyncClose()

		go func(partitionConsumer sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				fmt.Println(message.Partition, "", string(message.Key), " ", string(message.Value))
				var ld = new(es.LogData)
				if err := json.Unmarshal(message.Value, &ld.Data); err != nil {
					slog.Error("json unmarshal error ", "Error", err)
					continue
				}
				ld.Topic = message.Topic
				es.SendToChan(ld)
			}
		}(pc)
	}
	select {}
}
