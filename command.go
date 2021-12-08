package main

import (
	"github.com/urfave/cli"
	"kafka-tools/copy"
	"kafka-tools/publish"
)

var copyCommand = cli.Command{
	Name:   "copy",
	Usage:  "Copy messages from source topic to the target topic",
	Action: copy.CommandHandler,
	Flags:  copy.FLAGS,
}

var publishCommand = cli.Command{
	Name:   "publish",
	Usage:  "Publish message to kafka topic",
	Action: publish.CommandHandler,
	Flags:  publish.FLAGS,
}

var Commands = []cli.Command{copyCommand, publishCommand}
