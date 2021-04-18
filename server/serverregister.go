package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Concentrate/go_module_test/model"
	"github.com/Concentrate/go_module_test/utils"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
	"time"
)

type ServerRegister interface {
	RegisterServer(info *model.ServerInstanceInfo) (bool, error)
}

type serverRegisterHandler struct {
	zkConnection *zk.Conn
}

var registerListenerMap = make(map[int]string)

// 注册自身的服务地址即可
func (serverRegisterHandler *serverRegisterHandler) RegisterServer(info *model.ServerInstanceInfo) (bool, error) {
	if info == nil || utils.IsStringEmpty(info.ServiceName) {
		return false, errors.New("invailed register, address or name can't be empty")
	}
	if info.ListenPort <= 1024 {
		return false, errors.New("please listener , port should greater than 1024")
	}
	selfIp, err := utils.GetSelfIp()
	if err != nil {
		return false, errors.New("cannot get local ip")
	}
	info.Address = selfIp + ":" + strconv.Itoa(info.ListenPort)
	registerListenerMap[info.ListenPort] = info.ServiceName

	if strings.HasSuffix(info.ServiceName, "/") {
		return false, errors.New("serviceName can't be endwith /")
	}
	jsonObj, _ := json.Marshal(info)
	encodingVa := base64.StdEncoding.EncodeToString(jsonObj)

	var registerParentPath = utils.GetRpcCommonRegistrPath(info.ServiceName)
	var subChildNode = registerParentPath + "/" + encodingVa
	tmpErr := createPath(serverRegisterHandler.zkConnection, subChildNode)
	if tmpErr != nil {
		return false, tmpErr
	}
	return true, nil
}

func createPath(con *zk.Conn, subChildNode string) error {
	curTime, _ := time.Now().MarshalText()
	var acls = zk.WorldACL(zk.PermAll) //控制访问权限模式
	var isExist bool
	var tmpError error
	_, tmpError = con.Create(subChildNode, curTime, zk.FlagEphemeral, acls)
	if tmpError == nil {
		return nil
	}
	if tmpError == zk.ErrNoNode {
		var tmpSubChildNode = subChildNode
		var lastSlash = strings.LastIndex(tmpSubChildNode, "/")
		tmpSubChildNode = tmpSubChildNode[:lastSlash]
		var lastNotExSuccPath = make([]string, 3)
		for {
			isExist, _, tmpError = con.Exists(tmpSubChildNode)
			if isExist && tmpError == nil {
				break
			}
			_, tmpError = con.Create(tmpSubChildNode, curTime, zk.PermCreate, acls)
			if tmpError == nil {
				for i := len(lastNotExSuccPath) - 1; i >= 0; i-- {
					_, _ = con.Create(lastNotExSuccPath[i], curTime, zk.PermCreate, acls)
				}
				break
			} else {
				lastNotExSuccPath = append(lastNotExSuccPath, tmpSubChildNode)
				lastSlash := strings.LastIndex(tmpSubChildNode, "/")
				if lastSlash == -1 {
					break
				}
				tmpSubChildNode = tmpSubChildNode[:lastSlash]
				if utils.IsStringEmpty(tmpSubChildNode) {
					break
				}
			}
		}
	}
	_, tmpError = con.Create(subChildNode, curTime, zk.FlagEphemeral, acls)
	return tmpError
}

// param is zkServerAddress
func NewRpcServerRegister(zkServerAddress string) (ServerRegister, error) {
	if utils.IsStringEmpty(zkServerAddress) {
		return nil, errors.New("zkServerAddress cannot be empty")
	}

	var hosts = []string{zkServerAddress}
	zkConn, _, err := zk.Connect(hosts, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return &serverRegisterHandler{
		zkConnection: zkConn,
	}, nil
}
