package mock

import (
	"os"
	"bufio"
	//"encoding/hex"
	"fmt"
	"strconv"
	//"strings"
	"time"
	sender "bluesky-protocol/agentMock/mock/sender"
	"bluesky-protocol/agentMock/mock/config"
	"encoding/hex"
	"strings"
)

func open() {
	//msgFile := "./agentMock/mock/msg.txt"
	//fin, err := os.Open(msgFile)
	//
	//if err != nil {
	//	fmt.Println("open file error!", err)
	//	return
	//}
	//defer fin.Close()
	//
	//rd := bufio.NewReader(fin)
	//for {
	//	line, isPrefix,err := rd.ReadLine()
	//	fmt.Println(isPrefix,err)
	//	for err == nil && !isPrefix {
	//		sender.Send(line)
	//		time.Sleep(time.Duration(config.Config().Interval) * time.Second)
	//	}
	//	if isPrefix {
	//		logger.Info.Println("buffer size to small")
	//		break
	//	}
	//	if err != io.EOF {
	//		logger.Info.Println(err)
	//		break
	//	} else {
	//		logger.Info.Println("读取文件完成")
	//		break
	//	}
	//}
	msgFile := "./agentMock/mock/msg.txt"
	fin, err := os.Open(msgFile)
	if err != nil {
		fmt.Println("open file error!", err)
		return
	}
	defer fin.Close()

	rd := bufio.NewReader(fin)
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
		sender.Send(binMsg)
		time.Sleep(time.Duration(config.Config().Interval) * time.Second)
	}
}

func HexToBye(hex string) []byte {
	length := len(hex) / 2
	slice := make([]byte, length)
	rs := []rune(hex)

	for i := 0; i < length; i++ {
		s := string(rs[i*2: i*2+2])
		value, _ := strconv.ParseInt(s, 16, 10)
		slice[i] = byte(value & 0xFF)
	}
	return slice
}

func Start() {
	open()
}
