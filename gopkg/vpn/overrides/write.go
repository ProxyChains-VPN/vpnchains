package overrides

import (
	"syscall"
	"vpnchains/gopkg/vpn"
)

func Write(fd int, msg []byte) (n int, err error) {
	if vpn.tcpConnsMap[fd] == nil {
		return syscall.Write(fd, msg)
	}

	return (*vpn.tcpConnsMap[fd]).Write(msg)
}
