package wireguard

// todo split into files

import (
	"log"
	"strconv"
	"syscall"
	"time"
)

func (tunnel *WireguardTunnel) connect4(fd int, sa *syscall.SockaddrInet4) (err error) {
	address := strconv.Itoa(int(sa.Addr[0])) + "." +
		strconv.Itoa(int(sa.Addr[1])) + "." +
		strconv.Itoa(int(sa.Addr[2])) + "." +
		strconv.Itoa(int(sa.Addr[3])) + ":" +
		strconv.Itoa(sa.Port)

	socket, err := tunnel.Net.Dial("tcp", address) // todo
	if err != nil {
		return err
	}

	tunnel.TcpFdMap[fd] = &socket
	return nil
}

func (tunnel *WireguardTunnel) Connect(fd int, sa syscall.Sockaddr) (err error) {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return tunnel.connect4(fd, sa)
	case *syscall.SockaddrInet6:
		return nil // todo tmp
	case *syscall.SockaddrUnix:
		return nil // todo кинуть ошибку
	}
	return nil
}

func (tunnel *WireguardTunnel) Read(fd int, buf []byte) (n int, err error) {
	if tunnel.TcpFdMap[fd] == nil {
		log.Println("fd not found, not tcp")
		return 0, nil
	}

	socket := tunnel.TcpFdMap[fd]
	if err := (*socket).SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil { // todo зачем??
		return -1, err
	}
	return (*socket).Read(buf)
}

func (tunnel *WireguardTunnel) Write(fd int, buf []byte) (n int, err error) {
	if tunnel.TcpFdMap[fd] == nil {
		return syscall.Write(fd, buf)
	}

	return (*tunnel.TcpFdMap[fd]).Write(buf)
}
