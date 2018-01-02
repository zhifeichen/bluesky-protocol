package sender

import (
	"errors"
	"fmt"
	"net"
	"github.com/zhifeichen/bluesky-protocol/agentMock/mock/config"
	"time"
)

func Send(msg []byte) error {
	fmt.Println("start send...")
	config := config.Config()
	addr := config.ServerAddr
	fmt.Println("addr: ", addr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("resolve addr error: ", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("dial tcp error: ", err, tcpAddr)
		return err
	}
	defer conn.Close()

	split := len(msg)/2

	ret, err := conn.Write(msg[:split])
	fmt.Println("send data:",msg[:split])
	time.Sleep(time.Duration(200)* time.Millisecond)
	ret1, err:= conn.Write(msg[split:])
	fmt.Println("send data:",msg[split:])
	ret += ret1
	if err != nil {
		fmt.Println("write msg error: ", err)
		return err
	}
	if ret != len(msg) {
		fmt.Println("ret not equ send error: ", err)
		return errors.New("send error!")
	}
	fmt.Println("发送数据:",msg,"| ret:",ret,"...  [ok]")
	return nil
}