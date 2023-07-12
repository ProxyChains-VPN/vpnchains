package wireguard

// todo split into files

import (
	"log"
	"strconv"
	"syscall"
	"time"
)

func (tunnel *WireguardTunnel) connect4(fd int32, sa *syscall.SockaddrInet4) (err error) {
	address := strconv.Itoa(int(sa.Addr[0])) + "." +
		strconv.Itoa(int(sa.Addr[1])) + "." +
		strconv.Itoa(int(sa.Addr[2])) + "." +
		strconv.Itoa(int(sa.Addr[3])) + ":" +
		strconv.Itoa(sa.Port)

	log.Println(address)

	//socket, err := tunnel.Net.Dial("tcp4", address) // todo
	socket, err := tunnel.Net.Dial("tcp4", address)
	if err != nil {
		log.Println(err, "24 line overrides")
		return err
	}

	tunnel.TcpFdMap[fd] = &socket
	return nil
}

func (tunnel *WireguardTunnel) Connect(fd int32, sa syscall.Sockaddr) (err error) {
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

func (tunnel *WireguardTunnel) Read(fd int32, buf []byte) (n int64, err error) {
	if tunnel.TcpFdMap[fd] == nil {
		log.Println("fd not found, not tcp")
		return 0, nil
	}

	socket := tunnel.TcpFdMap[fd]
	if err := (*socket).SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil { // todo зачем??
		return -1, err
	}

	res, err := (*socket).Read(buf)
	return int64(res), err
}

func (tunnel *WireguardTunnel) Write(fd int32, buf []byte) (n int64, err error) {
	if tunnel.TcpFdMap[fd] == nil {
		log.Println("fd not found, not tcp")
		return 0, nil
		//res, err := syscall.Write(int(fd), buf)
		//return int64(res), err
	}

	res, err := (*tunnel.TcpFdMap[fd]).Write(buf)
	return int64(res), err
}
