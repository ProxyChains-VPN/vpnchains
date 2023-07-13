package vpn

import (
	"net"
	"syscall"
)

type Tunnel interface {
	Connect(fd int32, sa syscall.Sockaddr) (net.Conn, error)
}
