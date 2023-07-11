package overrides

import (
	"syscall"
	"vpnchains/gopkg/vpn"
)

func Close(fd int) (err error) {
	socket := vpn.tcpConnsMap[fd]
	if socket == nil {
		return syscall.Close(fd)
	}

	vpn.tcpConnsMap[fd] = nil
	return (*socket).Close()
}
