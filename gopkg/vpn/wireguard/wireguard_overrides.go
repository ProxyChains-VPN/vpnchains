package wireguard

// todo split into files

import (
	"gvisor.dev/gvisor/pkg/errors"
	"log"
	"net"
	"strconv"
	"syscall"
)

func (tunnel *WireguardTunnel) connect4(fd int32, sa *syscall.SockaddrInet4) (err error) {
	address := strconv.Itoa(int(sa.Addr[0])) + "." +
		strconv.Itoa(int(sa.Addr[1])) + "." +
		strconv.Itoa(int(sa.Addr[2])) + "." +
		strconv.Itoa(int(sa.Addr[3])) + ":" +
		strconv.Itoa(sa.Port) // todo будто бы можно без этого обойтись

	log.Println(address)

	//socket, err := tunnel.Net.Dial("tcp4", address) // todo
	socket, err := tunnel.Net.Dial("tcp4", address)
	if err != nil {
		log.Println(err, "24 line overrides")
		return err
	}

	go func(socket net.Conn) {
		buf := make([]byte, 32768)
		for { // TodO сделать норм мультиплексирование
			n, err := socket.Read(buf)
			if err != nil {
				log.Println(err, "31 line overrides")
				return
			}
			log.Println("read from socket", n)
			_, err = socket.Write(buf[:n])
			if err != nil {
				log.Println(err, "37 line overrides")
			}
		}
	}(socket)

	return nil
}

func (tunnel *WireguardTunnel) Connect(fd int32, sa syscall.Sockaddr) (err error) {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return tunnel.connect4(fd, sa)
	case *syscall.SockaddrInet6:
		return errors.New(0, "ipv6 not supported")
	case *syscall.SockaddrUnix:
		return errors.New(0, "unix sockets are not supposed to be here")
	}
	return nil
}
