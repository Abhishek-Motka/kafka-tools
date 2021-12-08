package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main() {
	cliApp := &cli.App{
		Name:                   "kafka-tools",
		Version:                "1.0.0",
		Description:            "Helper utilities to work with the kafka setup",
		Compiled:               time.Time{},
		Authors:                []cli.Author{{"Abhishek Patel", "abhishek.patel@gojek.com",}},
		UseShortOptionHandling: true,
		Commands: 				Commands,
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		fmt.Errorf("Failed to run command: %v", err)
	}
}
