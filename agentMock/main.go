package main

import (
	"bluesky-protocol/agentMock/mock/msg"
	"bluesky-protocol/agentMock/mock/config"
)

func main() {
	config.ParseConfig("./config.json")
	mock.Start()
}
