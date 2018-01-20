package common

import "errors"

var (
	NOT_IMPLEMENT_ERROR = errors.New("方法未实现")

)

type BaseError struct {
	code int
	msg  string
}

func (e BaseError) Error() string {
	return e.msg
}

func NewError(code int) *BaseError {
	return &BaseError{
		code,
		GetMsgByCode(code),
	}
}

func NewErrorOfMsg(code int, msg string) *BaseError {
	return &BaseError{
		code,
		msg,
	}
}
