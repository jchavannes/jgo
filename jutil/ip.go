package jutil

import (
	"fmt"
	"net"
)

func GetIpString(ip net.IP, port uint16) string {
	if len(ip) > 4 {
		return fmt.Sprintf("[%s]:%d", ip.String(), port)
	} else {
		return fmt.Sprintf("%s:%d", ip.String(), port)
	}
}
