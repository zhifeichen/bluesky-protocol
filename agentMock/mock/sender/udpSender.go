package sender

import (
	"time"
	"strings"
	"encoding/hex"
	"bufio"
	"errors"
	"fmt"
	"github.com/zhifeichen/bluesky-protocol/agentMock/mock/config"
	"net"
)

func UDPSend(msg []byte) error {
	fmt.Println("start send...")
	config := config.Config()
	addr := config.ServerAddr
	fmt.Println("addr: ", addr)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("resolve addr error: ", err)
		return err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("dial tcp error: ", err, udpAddr)
		return err
	}
	defer conn.Close()

	ret, err := conn.Write(msg)
	fmt.Println("send data:", msg)
	if err != nil {
		fmt.Println("write msg error: ", err)
		return err
	}
	if ret != len(msg) {
		fmt.Println("ret not equ send error: ", err)
		return errors.New("send error")
	}
	fmt.Println("发送数据:", msg, "| ret:", ret, "...  [ok]")
	return nil
}

func UDPSendFile(rd *bufio.Reader) error {
	fmt.Println("start send...")
	config := config.Config()
	addr := config.ServerAddr
	fmt.Println("addr: ", addr)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("resolve addr error: ", err)
		return err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("dial tcp error: ", err, udpAddr)
		return err
	}
	defer conn.Close()

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
		conn.Write(binMsg)
		// sender.UDPSend(binMsg)
		time.Sleep(time.Duration(config.Interval) * 100 * time.Millisecond)
	}
	return nil
}
