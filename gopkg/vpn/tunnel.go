package vpn

import "syscall"

type Tunnel interface {
	Connect(fd int, sa syscall.Sockaddr) (err error)
	Read(fd int, buf []byte) (n int, err error)
	Write(fd int, buf []byte) (n int, err error)
}
