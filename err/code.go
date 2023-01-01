package err

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
	return e.Msg
}

//func (e *Error) Code() string {
//	return fmt.Sprintf()
//}

func (e *Error) WithData(data []byte) *Error {
	return &Error{
		Code: e.Code,
		Msg:  e.Msg,
		Data: data,
	}
}

func FromError(err error) (s *Error, ok bool) {
	if err == nil {
		return nil, true
	}
	if se, ok := err.(*Error); ok {
		return se, ok
	}

	return nil, false
}
