package vpn

import (
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"net"
)

type Tunnel interface {
	Connect(addr *net.TCPAddr) (*gonet.TCPConn, error)
}
