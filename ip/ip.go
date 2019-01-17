package ip

import (
	"fmt"
	"net"
	"strings"
)

// Get public ip from 8.8.8.8
func GetPublicIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	ip := net.ParseIP(localAddr[0:idx])
	if ip == nil {
		ip = InputIP()
	}
	if !IsPublicIP(ip) {
		ip = InputIP()
	}
	return ip.String()
}

// Check whether the ip is public
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

// Get masternode ip address through user input
func InputIP() net.IP {
	var ip_input string
	fmt.Println("Input your public IP Address: ")
ERR:
	fmt.Scanln(&ip_input)
	ip := net.ParseIP(ip_input)
	if ip == nil {
		fmt.Println("Wrong IP Address, Input again: ")
		goto ERR
	}
	return ip
}
