package ip

import (
	"errors"
	"github.com/xjh22222228/ip"
	"net"
)

// local ip
func LocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for index := range addrs {
		// 检查ip地址判断是否回环地址
		if IPNet, ok := addrs[index].(*net.IPNet); ok && !IPNet.IP.IsLoopback() {
			if IPNet.IP.To4() != nil {
				return IPNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("failed to found IP address")
}

// ip v4
func IPV4() (string, error) {
	return ip.V4()
}

// ip v6
func IPV6() (string, error) {
	return ip.V6()
}

func IsIPv4(s string) bool {
	return ip.IsIPv4(s)
}

func IsIPv6(s string) bool {
	return ip.IsIPv6(s)
}
