package wireguard

import (
	"context"
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"log"
	"net"
	"time"
)

// Connect connects to the specified address.
// Is used when intercepting connect() syscalls.
// As sometimes connect() is called on non-blocking sockets, and we have blocking ones,
// timeout is set to 4 seconds, so browsers are not stuck.
// addr - address to connect to.
func (tunnel *WireguardTunnel) Connect(addr *net.TCPAddr) (*gonet.TCPConn, error) { // todo достаточно net.Conn????
	dialContext, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*4))
	defer cancel()

	socket, err := tunnel.net.DialContextTCP(dialContext, addr)
	if err != nil {
		log.Println(err, "24 line overrides")
		return nil, err
	}

	return socket, nil
}
