package overrides

import (
	"syscall"
	"time"
	"vpnchains/gopkg/vpn"
)

func Read(fd int, buf []byte) (n int, err error) {
	if vpn.tcpConnsMap[fd] == nil {
		return syscall.Read(fd, buf)
	}

	socket := vpn.tcpConnsMap[fd]
	if err := (*socket).SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		return -1, err
	}
	return (*socket).Read(buf)
}
