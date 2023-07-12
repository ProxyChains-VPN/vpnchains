package vpn

import (
	"syscall"
)

type Tunnel interface {
	Connect(fd int32, sa syscall.Sockaddr) (err error)
	Read(fd int32, buf []byte) (n int64, err error)
	Write(fd int32, buf []byte) (n int64, err error)
	Close(fd int32) (err error)
}
