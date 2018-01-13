package server

import (
	"bufio"
	"net"
)

type Handler interface {
	Handle(interface{}, WriteCloser) error
	HandleUdp(interface{}, WriteCloser,*net.UDPAddr) error
}

type Codec interface {
	GetScanSplitFun() bufio.SplitFunc		// call by read
	Decode([]byte) (interface{}, error)		// call by read
	Encode(interface{}) ([]byte, error)		// call by write
}

// ContextKey is the key type for putting context-related data.
type contextKey string

// Context keys for messge, server and net ID.
const (
	messageCtx contextKey = "message"
	serverCtx  contextKey = "server"
	netIDCtx   contextKey = "netid"
)
