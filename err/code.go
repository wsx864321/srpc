package err

import "fmt"

const (
	Ok = 0

	UnKnowErr           = 1000
	ServerErr           = 1001
	ServiceNotExistCode = 1002
	MethodNotExistCode  = 1003
)

var (
	OkErr              = NewError(Ok, "")
	ServiceNotExistErr = NewError(ServiceNotExistCode, "server not exist")
	MethodNotExistErr  = NewError(MethodNotExistCode, "method not exist")
)

func NewError(code int32, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("type : server, code : %d, msg : %s", e.Code, e.Msg)
}

func (e *Error) WithData(data []byte) *Error {
	return &Error{
		Code: e.Code,
		Msg:  e.Msg,
		Data: data,
	}
}
