package utils

import (
	"github.com/samuel/go-zookeeper/zk"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

func IsStringEmpty(input string) bool {
	return len(input) <= 0
}

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func FlatCookiesToMap(cookies []*http.Cookie) map[string]string {
	if cookies == nil {
		return nil
	}
	cookieMapValue := make(map[string]string, len(cookies))
	for i := 0; i < len(cookies); i++ {
		tmpCookie := cookies[i]
		cookieMapValue[tmpCookie.Name] = tmpCookie.Value
	}
	return cookieMapValue
}

func createLastNodeEphemeral(con *zk.Conn, subChildNode string) error {
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
				if IsStringEmpty(tmpSubChildNode) {
					break
				}
			}
		}
	}
	_, tmpError = con.Create(subChildNode, curTime, zk.FlagEphemeral, acls)
	return tmpError
}
