package kafka

import (
	"github.com/IBM/sarama"
	"testing"
)

func TestKafka(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			t.Fatal(err)
		}
		defer pc.AsyncClose()

		go func(partitionConsumer sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				t.Log(message.Partition, "", string(message.Key), " ", string(message.Value))

			}
		}(pc)
	}
	select {}
}
