package client

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Listener interface {
	EventHandler(message *kafka.Message) (bool, error)
	ErrorHandler(error kafka.Error)
}
