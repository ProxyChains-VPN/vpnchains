package vpn

import (
	"syscall"
	"time"
)

func Read(fd int, buf []byte) (n int, err error) {
	if conns[fd] == nil {
		return syscall.Read(fd, buf)
	}

	socket := conns[fd]
	if err := (*socket).SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		return -1, err
	}
	return (*socket).Read(buf)
}
