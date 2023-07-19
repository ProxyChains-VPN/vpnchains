package wireguard

// todo split into files

import (
	"context"
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"log"
	"net"
	"time"
)

func (tunnel *WireguardTunnel) Connect(addr *net.TCPAddr) (*gonet.TCPConn, error) { // todo достаточно net.Conn????
	dialContext, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*4))
	defer cancel()

	socket, err := tunnel.Net.DialContextTCP(dialContext, addr)
	//socket, err := tunnel.Net.DialTCP(addr)
	//socket, err := net.DialTCP("tcp4", address)
	if err != nil {
		log.Println(err, "24 line overrides")
		return nil, err
	}

	return socket, nil
}
