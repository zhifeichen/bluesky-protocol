package mock

import (
	"bufio"
	"os"
	//"encoding/hex"
	"fmt"
	"strconv"
	//"strings"
	"encoding/hex"
	"github.com/zhifeichen/bluesky-protocol/agentMock/mock/config"
	sender "github.com/zhifeichen/bluesky-protocol/agentMock/mock/sender"
	"strings"
	"time"
)

func open() {
	msgFile := "./mock/msg.txt"
	fin, err := os.Open(msgFile)
	if err != nil {
		fmt.Println("open file error!", err)
		return
	}
	defer fin.Close()

	rd := bufio.NewReader(fin)

	sender.SendFile(rd)
	sender.UDPSendFile(rd)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			fmt.Println("read string error", err)
			if len(line) == 0 {
				break
			}
		}
		// binMsg := HexToBye(line)
		binMsg, err := hex.DecodeString(strings.TrimRight(line, "\n"))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s: %v\n", line, binMsg)
		// sender.Send(binMsg)
		// sender.UDPSend(binMsg)
		time.Sleep(time.Duration(config.Config().Interval) * time.Second)
	}
}

func open2() {
	msgFile := "./mock/msg.txt"
	fin, err := os.Open(msgFile)
	if err != nil {
		fmt.Println("open file error!", err)
		return
	}
	defer fin.Close()

	rd := bufio.NewReader(fin)

	// sender.SendFile(rd)
	sender.UDPSendFile(rd)
}

func HexToBye(hex string) []byte {
	length := len(hex) / 2
	slice := make([]byte, length)
	rs := []rune(hex)

	for i := 0; i < length; i++ {
		s := string(rs[i*2 : i*2+2])
		value, _ := strconv.ParseInt(s, 16, 10)
		slice[i] = byte(value & 0xFF)
	}
	return slice
}

func Start() {
	go open()
	open2()
}
