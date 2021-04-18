package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Concentrate/go_module_test/model"
	"github.com/Concentrate/go_module_test/utils"
	"io"
	"log"
	"net"
	"strconv"
)

type RpcHandler interface {
	// HandleRequest param and return are both struct
	HandleRequest(param interface{}) (interface{}, error)
}

func startToListenPort(port int) {
	log.Default().Println("server start to listener ", port)
	conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Errorf("listen port server open error ,%v", port)
		return
	}

	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Errorf("连接出错 ,port is %v", port)
		}
		go serverConnHandler(conn)
	}
}
func writeToRpcRemote(writer *bufio.ReadWriter, data []byte) {
	log.Default().Println("response data is ", string(data))
	writer.Write(data)
	writer.WriteByte(utils.END_OF_RPC_INPUT)
	writer.Flush()
}

func serverConnHandler(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Errorf("recover...: %v", r)
		}
	}()

	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	byteContent, err := readWriter.ReadBytes(utils.END_OF_RPC_INPUT)
	if err != nil && err != io.EOF {
		fmt.Errorf("server read data error ,%v", err)
		return
	}
	//log.Default().Println("raw request input is ", string(byteContent))

	byteContent = byteContent[:len(byteContent)-1]

	//log.Default().Println("after remove delimiter request input is ", string(byteContent))

	var request model.RPCRequest
	var data []byte
	err = json.Unmarshal(byteContent, &request)
	if err != nil {
		fmt.Errorf("json parse rpcrequest error ")
		return
	}
	log.Default().Println("parse  request is ", request)
	handler, exist := handlerMap[connectServiceMethod(request.ServiceName, request.MethodName)]
	var response = model.RPCResponse{}
	response.RequestId = request.RequestId

	if !exist {
		response.ErrorInfo = (&NoMethodError{}).Error()
		data, _ = json.Marshal(&response)
		writeToRpcRemote(readWriter, data)
		return
	}

	result, tmpErr := handler.HandleRequest(request.Params)
	if tmpErr != nil {
		response.ErrorInfo = tmpErr.Error()
		data, _ = json.Marshal(&response)
		writeToRpcRemote(readWriter, data)
		return
	}
	response.Result = result
	data, _ = json.Marshal(&response)
	writeToRpcRemote(readWriter, data)
	return
}

var handlerMap = map[string]RpcHandler{}

func connectServiceMethod(service string, method string) string {
	return service + "___" + method
}

func RegisterHanlder(serviceName string, methodName string, hanlder RpcHandler) (bool, error) {
	if utils.IsStringEmpty(methodName) || utils.IsStringEmpty(serviceName) || hanlder == nil {
		return false, errors.New("params should not be empty")
	}
	handlerMap[connectServiceMethod(serviceName, methodName)] = hanlder
	return true, nil
}

func StartListenRpcServer() {
	log.Default().Println("start rpc server listen")
	for k, v := range registerListenerMap {
		log.Default().Println("start to listen port", k, " and the provider service is ", v)
		go startToListenPort(k)
	}
}
