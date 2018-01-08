package tcpServer

import (
	"bufio"
)

type Handler interface {
	Handle(interface{}, WriteCloser) error
}

// HandlerFunc serves as an adapter to allow the use of ordinary functions as handlers.
type HandlerFunc func(interface{}, WriteCloser) error

// Handle calls f(ctx, c)
func (f HandlerFunc) Handle(msg interface{}, c WriteCloser) error{
	return f(msg, c)
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
