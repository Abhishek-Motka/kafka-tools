package publish

import "github.com/urfave/cli"

var (
	serverFlag = &cli.StringFlag{
		Name:     "server",
		Usage:    "Bootstrap server to publish messages",
		Required: false,
		Value:    "localhost:9092",
	}
	topicFlag = &cli.StringFlag{
		Name:     "topic",
		Usage:    "Topic on target servers to copy message to",
		Required: true,
	}
	messageFlag = &cli.StringFlag{
		Name:     "message",
		Usage:    "Number of messages to copy from source topic to target topic",
		Required: false,
	}
	protoFlag = &cli.StringFlag{
		Name:     "proto",
		Usage:    "Protobuf schema to encode message using protobuf",
		Required: false,
	}
	protoDirFlag = &cli.StringFlag{
		Name:     "proto-dir",
		Usage:    "Root directory for proto packages",
		Required: false,
	}
	messageNameFlag = &cli.StringFlag{
		Name:     "message-name",
		Usage:    "Proto message name",
		Required: false,
	}
)

var FLAGS = []cli.Flag{serverFlag, topicFlag, messageFlag, protoFlag, protoDirFlag, messageNameFlag}

func Server(ctx *cli.Context) string {
	return ctx.String("server")
}

func Topic(ctx *cli.Context) string {
	return ctx.String("topic")
}

func Message(ctx *cli.Context) string {
	return ctx.String("message")
}

func Proto(ctx *cli.Context) string {
	return ctx.String("proto")
}

func ProtoDir(ctx *cli.Context) string {
	return ctx.String("proto-dir")
}

func ProtoMessageName(ctx *cli.Context) string {
	return ctx.String("message-name")
}

func ShouldEncodeUsingProto(ctx *cli.Context) bool {
	return ctx.IsSet("proto")
}

func AreProtoFlagsSet(ctx *cli.Context) bool {
	return ctx.IsSet("proto-dir") && ctx.IsSet("message-name")
}

func IsMessageProvided(ctx *cli.Context) bool {
	return ctx.IsSet("message")
}
