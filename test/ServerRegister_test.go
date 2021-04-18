package test

import (
	"fmt"
	"github.com/Concentrate/go_module_test/client"
	"github.com/Concentrate/go_module_test/model"
	"github.com/Concentrate/go_module_test/server"
	"github.com/Concentrate/go_module_test/utils"
	"log"
	"sync"
	"testing"
	"time"
)

const zkServer = "localhost:2181"

func TestHelloworld(t *testing.T) {
	fmt.Println("hello,world")
}

func TestServerReigst(t *testing.T) {
	reigster, err := server.NewRpcServerRegister(zkServer)
	if err != nil {
		log.Default().Fatalln(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(5)
	var reigestFun = func(serviceName string, port int) {
		if _, err = reigster.RegisterServer(&model.ServerInstanceInfo{ServiceName: serviceName,
			ListenPort: port,
		}); err != nil {
			log.Default().Println(err)
		}
	}
	reigestFun("HelloServiceTest1", 8761)
	reigestFun("HelloServiceTest1", 8861)
	reigestFun("HelloServiceTest1", 8862)
	reigestFun("HelloServiceTest1", 8863)

	reigestFun("HelloServiceTest2", 8762)
	reigestFun("HelloServiceTest2", 8861)

	reigestFun("HelloServiceTest3", 8763)
	reigestFun("HelloServiceTest3", 8862)

	server.StartListenRpcServer()

	time.Sleep(time.Hour)
}

func TestServerReigst2(t *testing.T) {
	reigster, err := server.NewRpcServerRegister(zkServer)
	if err != nil {
		log.Default().Fatalln(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(5)
	var reigestFun = func(serviceName string, port int) {
		if _, err = reigster.RegisterServer(&model.ServerInstanceInfo{ServiceName: serviceName,
			ListenPort: port,
		}); err != nil {
			log.Default().Println(err)
		}
	}
	reigestFun("HelloServiceTest1", 8885)

	server.StartListenRpcServer()

	time.Sleep(time.Hour)
}

func TestConsumer(t *testing.T) {
	client, _ := client.NewRpcClientRequest(zkServer)
	go client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest1", MethodName: "hello"})
	go client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest2", MethodName: "hello"})
	go client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest3", MethodName: "hello"})
	time.Sleep(10 * time.Second)
}

func TestUtils(t *testing.T) {
	log.Default().Println(utils.GetSelfIp())

}

type TmpParam struct {
	Name  string
	Local string
}

func TestClientCall(t *testing.T) {
	client, _ := client.NewRpcClientRequest(zkServer)
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest1", MethodName: "hello", Params: &TmpParam{
		Name: "ok", Local: "guangdong",
	}})
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest1", MethodName: "hello", Params: &TmpParam{
		Name: "ok", Local: "guangdong",
	}})
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest1", MethodName: "hello", Params: &TmpParam{
		Name: "ok", Local: "guangdong",
	}})
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest1", MethodName: "hello", Params: &TmpParam{
		Name: "ok", Local: "guangdong",
	}})
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest2", MethodName: "hello"})
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest2", MethodName: "hello"})
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest3", MethodName: "hello"})
	client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest3", MethodName: "hello"})

	time.Sleep(time.Hour)
}

func TestClientCall2(t *testing.T) {
	client, _ := client.NewRpcClientRequest(zkServer)

	for i := 0; i < 100; i++ {
		client.CallRpcRequest(&model.RPCRequest{ServiceName: "HelloServiceTest1", MethodName: "hello", Params: &TmpParam{
			Name: "ok", Local: "guangdong",
		}})
		time.Sleep(3 * time.Second)
	}

	time.Sleep(time.Hour)
}

func TestZkPathUtils(t *testing.T) {
	tmp1 := utils.GetRpcCommonRegistrPath("helloMyworld")
	log.Default().Println(utils.ParseServiceNameFromZkPath(tmp1))

	var tmpLen []string
	var tmpMap map[string]string
	log.Default().Println(tmpLen, len(tmpLen), tmpMap, len(tmpMap))
}
