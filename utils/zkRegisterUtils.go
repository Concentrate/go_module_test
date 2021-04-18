package utils

import "strings"

const common_PATH_PREFIX = "/rpc/"
const service_INNER_PREFIX = "/service"

func GetRpcCommonRegistrPath(serviceName string) string {
	var registerParentPath = common_PATH_PREFIX + serviceName + service_INNER_PREFIX
	return registerParentPath
}

func ParseServiceNameFromZkPath(zkPath string) string {

	index := strings.Index(zkPath, common_PATH_PREFIX)
	lastSubIndex := strings.LastIndex(zkPath, service_INNER_PREFIX)

	if index == -1 || lastSubIndex == -1 {
		return ""
	}
	var startIndex = index + len(common_PATH_PREFIX)
	if startIndex >= len(zkPath) || startIndex >= lastSubIndex {
		return ""
	}
	serviceName := zkPath[startIndex:lastSubIndex]
	return serviceName
}
