package wireguard

// todo split into files

import (
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"log"
	"net"
)

func (tunnel *WireguardTunnel) Connect(addr *net.TCPAddr) (*gonet.TCPConn, error) { // todo достаточно net.Conn????
	socket, err := tunnel.Net.DialTCP(addr)
	//socket, err := net.DialTCP("tcp4", address)
	if err != nil {
		log.Println(err, "24 line overrides")
		return nil, err
	}

	return socket, nil
}
