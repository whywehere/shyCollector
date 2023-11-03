package utils

import (
	"log/slog"
	"net"
	"strings"
	"testing"
)

func TestGETLocalIP(t *testing.T) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	loadAddr := conn.LocalAddr().(*net.UDPAddr)
	slog.Info(loadAddr.String())
	ip := strings.Split(loadAddr.IP.String(), ":")[0]
	t.Log(ip)
}
