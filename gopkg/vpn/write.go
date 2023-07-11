package vpn

import "syscall"

func Write(fd int, msg []byte) (n int, err error) {
	if conns[fd] == nil {
		return syscall.Write(fd, msg)
	}

	return (*conns[fd]).Write(msg)
}
