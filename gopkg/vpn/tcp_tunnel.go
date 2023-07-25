package vpn

import (
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"net"
)

// TcpTunnel An interface that represents a VPN tunnel.
// Connect(addr *net.TCPAddr) - connects to the specified address.
// Is used when intercepting connect() syscalls.
type TcpTunnel interface {
	Connect(addr *net.TCPAddr) (*gonet.TCPConn, error)
}
