package copy

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli"
	"kafka-tools/client"
	"sync/atomic"
	"time"
)

type copyCommandListner struct {
	count          int64
	ctx            *cli.Context
	messageChannel chan *kafka.Message
}

func (c *copyCommandListner) EventHandler(message *kafka.Message) (bool, error) {
	atomic.AddInt64(&c.count, 1)
	fmt.Printf("Message consumed from topic %s at offset %v\n", *message.TopicPartition.Topic, message.TopicPartition.Offset)
	if c.count > NumberOfMessagesToConsume(c.ctx) {
		close(c.messageChannel)
		return true, nil
	}
	c.messageChannel <- message
	return false, nil
}

func (c *copyCommandListner) ErrorHandler(error kafka.Error) {
	fmt.Printf("Encountered Error In Kafka: %v", error)
}

func CommandHandler(ctx *cli.Context) error {
	messageChannel := make(chan *kafka.Message, 10)
	errorChannel := make(chan error)
	go func() {
		err := client.KafkaProduce(getKafkaProducerConfig(ctx), TargetTopic(ctx), messageChannel)
		errorChannel <- err
	}()
	copyCommandHandler := copyCommandListner{ctx: ctx, messageChannel: messageChannel}
	go func() {
		timeout := time.Duration(Timeout(ctx))
		kafkaConfig := getKafkaConsumerConfig(ctx)
		err := client.KafkaConsume(kafkaConfig, SourceTopic(ctx), &copyCommandHandler, timeout)
		errorChannel <- err
	}()
	err := <-errorChannel
	close(errorChannel)
	if err != nil {
		fmt.Printf("Encountered Error: %v\n", err)
	} else {
		fmt.Println("Message transfer completed :)")
	}
	return err
}

func getKafkaProducerConfig(ctx *cli.Context) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": TargetServer(ctx),
	}
}

func getKafkaConsumerConfig(ctx *cli.Context) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":     SourceServer(ctx),
		"broker.address.family": "v4",
		"group.id":              ConsumerGroup(ctx),
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	}
}
