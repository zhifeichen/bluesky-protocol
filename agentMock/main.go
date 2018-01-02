package main

import (
	"bluesky-protocol/agentMock/mock/msg"
	mockConfig "bluesky-protocol/agentMock/mock/config"
)

func main() {
	mockConfig.ParseConfig("./agentMock/mock/config.json")
	mock.Start()
}