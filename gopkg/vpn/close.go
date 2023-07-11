package vpn

import "syscall"

func Close(fd int) (err error) {
	socket := conns[fd]
	if socket == nil {
		return syscall.Close(fd)
	}

	conns[fd] = nil
	return (*socket).Close()
}
