package model

type ServerInstanceInfo struct {
	ServiceName string
	ListenPort  int    // listen port
	Address     string // current server ip and port,address, don't need to fill
	protocol    string // serialze method
	weight      int    // not open now
}

func NewServerInstance(serviceName string, listenePort int) *ServerInstanceInfo {
	return &ServerInstanceInfo{ServiceName: serviceName, ListenPort: listenePort}
}

type RPCRequest struct {
	RequestId   string
	ServiceName string
	MethodName  string
	Params      interface{} // must be struct pointer?
}

type RPCResponse struct {
	RequestId string
	Result    interface{}
	ErrorInfo string
}
