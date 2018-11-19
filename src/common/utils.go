package common

// RequestPkt ...
type RequestPkt struct {
	MethodName string
	Args       []interface{}
}

// NewRequestPkt ...
func NewRequestPkt(methodName string, args ...interface{}) RequestPkt {
	var _args []interface{}
	for _, arg := range args {
		_args = append(_args, arg)
	}
	return RequestPkt{
		MethodName: methodName,
		Args:       _args,
	}
}

// ReturnPkt ...
type ReturnPkt struct {
	MethodName  string
	ReturnValue interface{}
	Err         error
}

// Service type encapsulates service information (name, ip, port) fot the naming server
type Service struct {
	Name string
	IP   string
	Port int32
}

// ConsultPkt ...
type ConsultPkt struct {
	ConsultType string
	Data        []byte
}

// RequestInfo ...
type RequestInfo struct {
	Name string
}
