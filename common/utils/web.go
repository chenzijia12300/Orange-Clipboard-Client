package utils

import (
	"fmt"
	"net"
)

const DefaultIp = "127.0.0.1"

// GetLocalIP 尝试获得本机内网IP（适用于单网卡）
func GetLocalIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get local address failure", err.Error())
		return DefaultIp
	}
	for _, addr := range address {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return DefaultIp
}
