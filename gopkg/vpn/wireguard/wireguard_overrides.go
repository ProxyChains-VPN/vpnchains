package wireguard

// todo split into files

import (
	"errors"
	"log"
	"strconv"
	"syscall"
	"time"
)

func (tun *WireguardTunnel) connect4(fd int32, sa *syscall.SockaddrInet4) (err error) {
	address := strconv.Itoa(int(sa.Addr[0])) + "." +
		strconv.Itoa(int(sa.Addr[1])) + "." +
		strconv.Itoa(int(sa.Addr[2])) + "." +
		strconv.Itoa(int(sa.Addr[3])) + ":" +
		strconv.Itoa(sa.Port)

	log.Println(address)

	//socket, err := tun.Net.Dial("tcp4", address) // todo
	socket, err := tun.Net.Dial("tcp4", address)
	if err != nil {
		log.Println(err, "24 line overrides")
		return err
	}

	tun.TcpFdMap[fd] = &socket
	return nil
}

func (tun *WireguardTunnel) Connect(fd int32, sa syscall.Sockaddr) (err error) {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return tun.connect4(fd, sa)
	case *syscall.SockaddrInet6:
		return nil // todo tmp
	case *syscall.SockaddrUnix:
		return nil // todo кинуть ошибку
	}
	return nil
}

func (tun *WireguardTunnel) Read(fd int32, buf []byte) (n int64, err error) {
	if tun.TcpFdMap[fd] == nil {
		return -1, errors.New("no such tcp socket")
	}

	socket := tun.TcpFdMap[fd]
	if err := (*socket).SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil { // todo зачем??
		return -1, err
	}

	res, err := (*socket).Read(buf)
	return int64(res), err
}

func (tun *WireguardTunnel) Write(fd int32, buf []byte) (n int64, err error) {
	if tun.TcpFdMap[fd] == nil {
		return -1, errors.New("no such tcp socket")
	}

	res, err := (*tun.TcpFdMap[fd]).Write(buf)
	return int64(res), err
}

func (tun *WireguardTunnel) Close(fd int32) (err error) {
	if tun.TcpFdMap[fd] == nil {
		return errors.New("no such socket")
	}

	tun.TcpFdMap[fd] = nil
	return nil
}
