package common

import (
	"fmt"
	"net"
	"strings"
)

const (
	CanNotGetIp = "get ip failed"
)

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8,8,8,8:80")
	if err != nil {
		return 
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
