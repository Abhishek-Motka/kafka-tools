package client

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func KafkaConsume(config *kafka.ConfigMap, topic string, listener Listener, timeout time.Duration) error {
	client, err := kafka.NewConsumer(config)
	if err != nil {
		return fmt.Errorf("failed to start consumer for the source server: %v", err)
	}
	defer client.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error)
	go func() {
		time.Sleep(time.Second * timeout)
		done <- errors.New("terminated due to timeout")
		close(done)
	}()

	err = client.Subscribe(topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe a topic %s: %v", topic, err)
	}

	for {
		select {
		case sig := <-sigChan:
			return fmt.Errorf("caught %v signal: terminating", sig)
		case timeoutErr := <-done:
			return timeoutErr
		default:
			event := client.Poll(100)
			switch e := event.(type) {
			case *kafka.Message:
				terminate, err := listener.EventHandler(e)
				if err != nil {
					return err
				}
				if terminate {
					return nil
				}
			case kafka.Error:
				listener.ErrorHandler(e)
				return fmt.Errorf("encountered error while consuming messages: %v", err)
			}
		}
	}
}

func KafkaProduce(config *kafka.ConfigMap, topic string, msgChan chan *kafka.Message) error {

	client, err := kafka.NewProducer(config)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create a producer: %v", err))
	}

	fmt.Printf("Created Producer %v\n", client)

	deliveryChan := make(chan kafka.Event)

	for {
		msg, channelOpen := <-msgChan
		fmt.Println(msg)
		if !channelOpen {
			close(deliveryChan)
			return nil
		}
		kafkaMsg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            msg.Key,
			Value:          msg.Value,
		}
		err = client.Produce(kafkaMsg, deliveryChan)
		if err != nil {
			close(deliveryChan)
			return errors.New(fmt.Sprintf("Failed to deliver message: %v", err))
		}
		event := <-deliveryChan
		m := event.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			close(deliveryChan)
			return errors.New(fmt.Sprintf("Failed to deliver message %v", m.TopicPartition.Error))
		} else {
			fmt.Printf("Delivered message to topic %s at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Offset)
		}
	}

	close(deliveryChan)
	return nil
}
