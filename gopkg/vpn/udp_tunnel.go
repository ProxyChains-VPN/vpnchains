package vpn

import (
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"net"
)

// UdpTunnel An interface that represents a VPN tunnel.
// Dial(network string, laddr, raddr *net.UDPAddr) - connects to the specified address.
// Is used when intercepting sendto() syscalls.
type UdpTunnel interface {
	Dial(addr *net.UDPAddr) (*gonet.UDPConn, error)
}
