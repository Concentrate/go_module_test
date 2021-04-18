package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Concentrate/go_module_test/model"
	"github.com/Concentrate/go_module_test/utils"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strconv"
	"time"
)

type ICallRcpRequest interface {
	CallRpcRequest(request *model.RPCRequest) (*model.RPCResponse, error)
}

type clientRpcCallHanlder struct {
	zkConnection *zk.Conn
}

func detailRpcRequestVailed(request *model.RPCRequest) bool {
	if request == nil {
		return false
	}
	if utils.IsStringEmpty(request.ServiceName) || utils.IsStringEmpty(request.MethodName) {
		return false
	}
	return true
}
func writeToRpcRemote(writer *bufio.ReadWriter, data []byte) (int, error) {
	n, err := writer.Write(data)
	writer.WriteByte(utils.END_OF_RPC_INPUT)
	writer.Flush()
	return n, err
}

func (receiver *clientRpcCallHanlder) CallRpcRequest(request *model.RPCRequest) (*model.RPCResponse, error) {
	if !detailRpcRequestVailed(request) {
		return nil, errors.New("request is not vailed, lack of serviceName or methodName")
	}
	address := SelectOneRemoteAddress(request.ServiceName)
	request.RequestId = utils.GetUUID()

	log.Default().Println(address)
	if utils.IsStringEmpty(address) {
		return nil, errors.New("no service provider right now")
	}
	connect, err := getTcpConnect(address)
	if err != nil {
		return nil, err
	}
	readWiter := bufio.NewReadWriter(bufio.NewReader(connect), bufio.NewWriter(connect))

	defer connect.Close()
	if err != nil {
		fmt.Errorf("%v", err)
		return nil, err
	}
	tmpResponse, err := requestToRemoteTcpServer(request, readWiter)
	if err != nil {
		log.Default().Println(err)
		for index := 0; innerCallConfig != nil && index < innerCallConfig.retryTimes && err != nil; index++ {
			log.Default().Println("retry strategy ,start to retry request remote and index is ", index, "retry times :",
				strconv.Itoa(innerCallConfig.retryTimes))
			tmpResponse, err = requestToRemoteTcpServer(request, readWiter)
		}

		if err != nil {
			return nil, err
		}
	}
	log.Default().Println("request service ", request,
		"  from ip address ", address, "tcp dialog response is ", tmpResponse)
	return tmpResponse, nil
}

func requestToRemoteTcpServer(request *model.RPCRequest, readWiter *bufio.ReadWriter) (*model.RPCResponse, error) {
	byteContent, _ := json.Marshal(request)
	_, tmpErr := writeToRpcRemote(readWiter, byteContent)
	if tmpErr != nil {
		return nil, tmpErr
	}
	var responseData []byte
	responseData, tmpErr = readWiter.ReadBytes(utils.END_OF_RPC_INPUT)
	if tmpErr != nil {
		return nil, tmpErr
	}

	//log.Default().Println("raw default response Data is ", string(responseData))
	responseData = responseData[:len(responseData)-1]

	//log.Default().Println("after remove delimter response data is ", string(responseData))
	var tmpResponse model.RPCResponse
	json.Unmarshal(responseData, &tmpResponse)
	return &tmpResponse, nil
}

var innerZkConn *zk.Conn

type ClientCallConfig struct {
	retryTimes int // 失败时候自动重试
}

var innerCallConfig *ClientCallConfig

func NewRpcClientRequest(zkServerAddress string, clientConfigs ...*ClientCallConfig) (ICallRcpRequest, error) {
	if utils.IsStringEmpty(zkServerAddress) {
		return nil, errors.New("zkServerAddress cannot be empty")
	}
	if len(clientConfigs) > 0 {
		innerCallConfig = clientConfigs[0]
	}
	var hosts = []string{zkServerAddress}
	zkConn, _, err := zk.Connect(hosts, 10*time.Second)
	if err != nil {
		return nil, err
	}
	innerZkConn = zkConn
	return &clientRpcCallHanlder{
		zkConnection: zkConn,
	}, nil
}
