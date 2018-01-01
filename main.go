package main

import (
	"github.com/zhifeichen/bluesky-protocol/receiver"
	"os"
	"github.com/zhifeichen/bluesky-protocol/config"
)

func main() {
	err := config.ParseConfig("./config.json")
	if err != nil {
		os.Exit(1)
	}
	receiver.Start()
}
