package vpn

import (
	"syscall"
)

type Tunnel interface {
	Connect(fd int32, sa syscall.Sockaddr) error
}
