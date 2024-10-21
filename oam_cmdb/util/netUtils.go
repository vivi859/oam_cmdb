package util

import (
	"net"
)

//返回本机局域网ip
func LocalIp() string {
	return QueryCacheFirst[string](GetPublicCache(), "local_ip", func() string {
		conn, err := net.Dial("udp", "8.8.8.8:8")
		if err != nil {
			addrList, err := net.InterfaceAddrs()
			if err == nil {
				for _, address := range addrList {
					if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
						if ipNet.IP.IsPrivate() {
							return ipNet.IP.String()
						}
					}
				}
			}
			return ""
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	})

}
