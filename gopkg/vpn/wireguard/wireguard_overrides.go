package wireguard

// todo split into files

import (
	"gvisor.dev/gvisor/pkg/errors"
	"log"
	"net"
	"strconv"
	"syscall"
)

func (tunnel *WireguardTunnel) connect4(fd int32, sa *syscall.SockaddrInet4) (net.Conn, error) {
	address := strconv.Itoa(int(sa.Addr[0])) + "." +
		strconv.Itoa(int(sa.Addr[1])) + "." +
		strconv.Itoa(int(sa.Addr[2])) + "." +
		strconv.Itoa(int(sa.Addr[3])) + ":" +
		strconv.Itoa(sa.Port) // todo будто бы можно без этого обойтись

	log.Println(address)

	//socket, err := tunnel.Net.Dial("tcp4", address)
	socket, err := net.Dial("tcp4", address)
	if err != nil {
		log.Println(err, "24 line overrides")
		return nil, err
	}

	return socket, nil
}

func (tunnel *WireguardTunnel) Connect(fd int32, sa syscall.Sockaddr) (net.Conn, error) {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return tunnel.connect4(fd, sa)
	case *syscall.SockaddrInet6:
		return nil, errors.New(0, "ipv6 not supported")
	case *syscall.SockaddrUnix:
		return nil, errors.New(0, "why unix is here???")
	}

	return nil, errors.New(0, "unknown sockaddr type") // todo errno
}
