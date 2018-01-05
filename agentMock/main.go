package main

import (
	mockConfig "github.com/zhifeichen/bluesky-protocol/agentMock/mock/config"
	"github.com/zhifeichen/bluesky-protocol/agentMock/mock/msg"
)

func main() {
	mockConfig.ParseConfig("./agentMock/mock/config.json")
	mock.Start()
}
