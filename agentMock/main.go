package main

import (
	"github.com/zhifeichen/bluesky-protocol/agentMock/mock/msg"
	mockConfig "github.com/zhifeichen/bluesky-protocol/agentMock/mock/config"
)

func main() {
	mockConfig.ParseConfig("./agentMock/mock/config.json")
	mock.Start()
}