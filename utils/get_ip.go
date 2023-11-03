package utils

import (
	"log/slog"
	"net"
	"strings"
)

func GetOutBoundIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	loadAddr := conn.LocalAddr().(*net.UDPAddr)
	slog.Info(loadAddr.String())
	ip = strings.Split(loadAddr.IP.String(), ":")[0]
	return
}
