package utils

import (
	"net"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, nil
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, nil
	}
	defer listen.Close()
	return listen.Addr().(*net.TCPAddr).Port, nil
}
