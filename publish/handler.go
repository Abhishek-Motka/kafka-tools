package publish

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/urfave/cli"
	"kafka-tools/client"
	"os"
	"strings"
)

func CommandHandler(ctx *cli.Context) error {
	message, err := getEncodedMessage(ctx)
	if err != nil {
		fmt.Printf("Encountered error: %v", err)
		return err
	}
	err = publishMessageToKafka(ctx, message)
	if err != nil {
		fmt.Printf("Encountered error: %v", err)
		return err
	}
	return nil
}

func publishMessageToKafka(ctx *cli.Context, message []byte) error {
	messageChan := make(chan *kafka.Message)
	errorChan := make(chan error)
	go func() {
		errorChan <- client.KafkaProduce(getKafkaProducerConfig(ctx), Topic(ctx), messageChan)
	}()
	messageChan <- &kafka.Message{Value: message}
	err := <-errorChan
	close(messageChan)
	return err
}

func getKafkaProducerConfig(ctx *cli.Context) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers": Server(ctx),
	}
}

func getMessage(ctx *cli.Context) string {
	if IsMessageProvided(ctx) {
		return Message(ctx)
	}
	scanner := bufio.NewScanner(os.Stdin)
	input := make([]string, 0)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	return strings.Join(input, "\n")
}

func getEncodedMessage(ctx *cli.Context) ([]byte, error) {
	message := getMessage(ctx)
	if ShouldEncodeUsingProto(ctx) {
		if false == AreProtoFlagsSet(ctx) {
			return nil, errors.New("proto-dir and message-name are required if proto option is provided")
		}
		return parseMessageToProto(ctx, message)
	}
	return []byte(message), nil
}

func parseMessageToProto(ctx *cli.Context, message string) ([]byte, error) {
	md, err := getMessageDescriptor(ProtoDir(ctx), Proto(ctx), ProtoMessageName(ctx))
	if err != nil {
		return nil, err
	}
	dynamicMessage := dynamic.NewMessage(md)
	err = jsonpb.UnmarshalString(message, dynamicMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON to proto: %v", err)
	}
	return dynamicMessage.Marshal()
}

func getFileDescriptorByPath(importPaths []string, protoFile string) (*desc.FileDescriptor, error) {
	parser := protoparse.Parser{
		ImportPaths:           importPaths,
		InferImportPaths:      len(importPaths) == 0,
		IncludeSourceCodeInfo: true,
	}
	fds, err := parser.ParseFiles(protoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proto file %s: %v", protoFile, err)
	}
	for _, value := range fds {
		if strings.HasSuffix(protoFile, value.GetName()) {
			return value, nil
		}
	}
	return nil, fmt.Errorf("failed to get file descriptor for given proto %s", protoFile)
}

func getMessageDescriptor(protoDir string, protoFile string, messageName string) (*desc.MessageDescriptor, error) {
	fd, err := getFileDescriptorByPath([]string{protoDir}, protoFile)
	if err != nil {
		return nil, err
	}
	md := fd.FindMessage(messageName)
	if md == nil {
		return nil, fmt.Errorf("colud not find message %s in proto file %s", messageName, protoFile)
	}
	return md, nil
}
