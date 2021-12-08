package copy

import "github.com/urfave/cli"

var (
	sourceServerFlag = &cli.StringFlag{
		Name:        "source-server",
		Usage:       "Source Bootstrap server to copy messages from",
		Required:    false,
		Value:       "localhost:9092",
	}
	targetServerFlag = &cli.StringFlag{
		Name: "target-server",
		Usage: "Target Bootstrap server to copy messages to",
		Required: false,
		Value: "localhost:9092",
	}
	sourceTopic = &cli.StringFlag{
		Name: "source-topic",
		Usage: "Topic from source servers to copy message from",
		Required: true,
	}
	targetTopic = &cli.StringFlag{
		Name: "target-topic",
		Usage: "Topic on target servers to copy message to",
		Required: true,
	}
	consumerGroup = &cli.StringFlag{
		Name: "group",
		Usage: "Consumer group for the source topic",
		Required: true,
	}
	count = &cli.IntFlag{
		Name:        "count",
		Usage:       "Number of messages to copy from source topic to target topic",
		Required:    false,
		Value:       100,
	}
	timeout = &cli.IntFlag{
		Name: "timeout",
		Usage: "Timeout in seconds to stop consuming",
		Required: false,
		Value: 10,
	}
)

var FLAGS = []cli.Flag{sourceServerFlag, targetServerFlag, sourceTopic, targetTopic, consumerGroup, count, timeout}

func SourceServer(ctx *cli.Context) string {
	return ctx.String("source-server")
}

func TargetServer(ctx *cli.Context) string {
	return ctx.String("target-server")
}

func SourceTopic(ctx *cli.Context) string {
	return ctx.String("source-topic")
}

func TargetTopic(ctx *cli.Context) string {
	return ctx.String("target-topic")
}

func ConsumerGroup(ctx *cli.Context) string {
	return ctx.String("group")
}

func NumberOfMessagesToConsume(ctx *cli.Context) int64 {
	return ctx.Int64("count")
}

func Timeout(ctx *cli.Context) int64 {
	return ctx.Int64("timeout")
}