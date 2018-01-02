package main

import (
	"github.com/zhifeichen/bluesky-protocol/mock/config"
	"github.com/zhifeichen/bluesky-protocol/mock/msg"
)

func main() {
	config.ParseConfig("./config.json")
	mock.Start()
}
