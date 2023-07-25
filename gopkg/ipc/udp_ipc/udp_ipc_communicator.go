package tcp_ipc

import (
	"bytes"
	"log"
	"net"
)

// UdpIpcConnection A struct that represents an IPC connection.
// addr - a net.TCPAddr instance.
type UdpIpcConnection struct {
	addr    *net.UDPAddr
	bufSize int
}

// UdpIpcCommunicator An interface that represents an IPC communicator.
// Read(handler func(conn *net.UDPConn)) - reads from the local socket.
type UdpIpcCommunicator interface {
	Read(handler func(srcAddr *net.UDPAddr, buf []byte)) error
}

// NewConnection creates a new UdpIpcConnection instance.
// udpAddr - a net.UDPAddr instance.
// bufSize - buffer size.
func NewConnection(udpAddr *net.UDPAddr, bufSize int) UdpIpcCommunicator {
	return &UdpIpcConnection{addr: udpAddr, bufSize: bufSize}
}

// NewConnectionFromIpPort creates a new UdpIpcConnection instance.
// ip - ip address.
// port - port.
// bufSize - buffer size.
func NewConnectionFromIpPort(ip net.IP, port int, bufSize int) UdpIpcCommunicator {
	return &UdpIpcConnection{
		addr: &net.UDPAddr{
			IP:   ip,
			Port: port,
			Zone: "",
		},
		bufSize: bufSize,
	}
}

// Read listens to the local socket.
// handler - a function that handles the connection.
func (ipcConnection *UdpIpcConnection) Read(handler func(*net.UDPAddr, []byte)) error {
	socket, err := net.ListenUDP("udp", ipcConnection.addr)
	if err != nil {
		return err
	}

	buf := make([]byte, ipcConnection.bufSize)
	for {
		n, srcAddr, err := socket.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		go handler(srcAddr, bytes.Clone(buf[:n]))
	}
}
