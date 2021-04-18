package client

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Concentrate/go_module_test/model"
	"github.com/Concentrate/go_module_test/utils"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

var serviceAddressMap = map[string][]string{}
var lock = sync.Mutex{}
var roubinMap = map[string]int32{}

func listenerServiceChildNodeChange(zkEventChannel <-chan zk.Event) {
	select {
	case tmpEvent := <-zkEventChannel:
		serviceName := utils.ParseServiceNameFromZkPath(tmpEvent.Path)
		if utils.IsStringEmpty(serviceName) {
			return
		}
		if tmpEvent.Type == zk.EventNodeChildrenChanged {
			log.Default().Println(tmpEvent, "going to remove service sub address list ", "  before is ", serviceName, ":",
				serviceAddressMap[serviceName])
			lock.Lock()
			defer lock.Unlock()
			serviceAddressMap[serviceName] = []string{}
			log.Default().Println("after remove service is ", serviceName, serviceAddressMap[serviceName])
		}
		break

	}
}

func getServiceAddress(serviceName string) []string {
	if len(serviceAddressMap[serviceName]) != 0 {
		return serviceAddressMap[serviceName]
	}
	lock.Lock()
	defer lock.Unlock()
	var serviceParentPath = utils.GetRpcCommonRegistrPath(serviceName)
	childData, _, event, err := innerZkConn.ChildrenW(serviceParentPath)
	if err != nil || len(childData) == 0 {
		return []string{}
	}
	addressList := make([]string, 0)
	for _, v := range childData {
		var serviceInstance model.ServerInstanceInfo
		decodeVa, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			continue
		}

		err = json.Unmarshal(decodeVa, &serviceInstance)
		if err != nil {
			continue
		}
		addressList = append(addressList, serviceInstance.Address)
	}
	serviceAddressMap[serviceName] = addressList
	go listenerServiceChildNodeChange(event)
	return addressList
}

func SelectOneRemoteAddress(serviceName string) string {
	addressList := getServiceAddress(serviceName)
	log.Default().Println(addressList)
	if len(addressList) <= 0 {
		return ""
	}
	var atomicAdd = int32(0)
	if _, exist := roubinMap[serviceName]; !exist {
		roubinMap[serviceName] = atomicAdd
	}
	atomicAdd = roubinMap[serviceName]
	atomicAdd = atomic.AddInt32(&atomicAdd, 1)
	roubinMap[serviceName] = atomicAdd
	var index = int(atomicAdd) % len(addressList)
	return addressList[index]
}

//var syncMap sync.Map

func getTcpConnect(address string) (net.Conn, error) {
	//va, exist := syncMap.Load(address)
	//if exist && va != nil {
	//	var tmpConnect = reflect.TypeOf(va).(net.Conn)
	//	return tmpConnect, nil
	//}

	connect, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	//syncMap.Store(address, connect)
	return connect, nil
}
