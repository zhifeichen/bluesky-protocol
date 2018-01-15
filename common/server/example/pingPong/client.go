package pingPong

import (
	"net"
	"fmt"
	"os"
	"bytes"
	"encoding/binary"
	"io"
)

func StartTcpClientSendMsg(ip string, port int, msgsLen int, done chan int) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	for i := 0; i < msgsLen; i++ {
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.LittleEndian, MAGIC_NUM)
		msg := fmt.Sprintf("%s-%d", "ping", i)
		body := []byte(msg)

		binary.Write(&buffer, binary.LittleEndian, (int16)(len(body)))
		binary.Write(&buffer, binary.LittleEndian, body)
		c.Write(buffer.Bytes())

		lengthBytes := make([]byte, 4)
		_, err := io.ReadFull(c, lengthBytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println("读取长度: ", lengthBytes)

		var msgLen int16
		err = binary.Read(bytes.NewReader(lengthBytes[2:4]), binary.LittleEndian, &msgLen)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bodyBuffer := make([]byte, msgLen)
		//fmt.Println("msglen:", msgLen)
		_, err = io.ReadFull(c, bodyBuffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ack:", string(bodyBuffer))
		ackMsg := fmt.Sprintf("pong-%s", msg)
		ack := string(bodyBuffer)

		fmt.Println("udp client ack:", ack)
		if ack != ackMsg {
			fmt.Println("tcp client 接收数据错误, except:", ackMsg, " but:", ack)
			os.Exit(1)
		}
		done <- 1
	}

	fmt.Println("关闭tcp 客户端")
}

func StartUdpClientSendMsg(ip string, port int, msgLen int, done chan int) {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("resolve addr error: ", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("dial tcp error: ", err, udpAddr)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("udp client:", conn.LocalAddr())
	for i := 0; i < msgLen; i++ {
		var buffer bytes.Buffer
		msg := fmt.Sprintf("%s-%d", "ping", i)
		body := []byte(msg)
		binary.Write(&buffer, binary.LittleEndian, MAGIC_NUM)
		binary.Write(&buffer, binary.LittleEndian, (int16)(len(body)))
		binary.Write(&buffer, binary.LittleEndian, body)
		conn.Write(buffer.Bytes())

		bufferRed := make([]byte, 1500)
		if n, _, err := conn.ReadFromUDP(bufferRed); err == nil {
			ackMsg := fmt.Sprintf("pong-%s", msg)
			ack := string(bufferRed[4:n])

			fmt.Println("udp client ack:", ack)
			if ack != ackMsg {
				fmt.Println("udp client 接收数据错误, except:", ackMsg, " but:", ack)
				os.Exit(1)
			}
			done <- 1
		} else {
			fmt.Println(err)
		}
	}

	fmt.Println("关闭 udp 客户端")
}
