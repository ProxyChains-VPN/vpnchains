package vpn

import "syscall"

type Tunnel interface {
	Connect(fd int, sa syscall.Sockaddr) (err error)
	Read(fd int, p []byte) (n int, err error)
	Write(fd int, p []byte) (n int, err error)
}
