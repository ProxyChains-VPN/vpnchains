package tcp_ipc

import (
	"log"
	"net"
)

// UdpIpcConnection A struct that represents an IPC connection.
// Addr - a net.TCPAddr instance.
type UdpIpcConnection struct {
	Addr *net.UDPAddr
}

// UdpIpcCommunicator An interface that represents an IPC communicator.
// New() - creates a new UdpIpcConnection instance.
// Listen(handler func(conn net.Conn)) - listens to the local socket.
type UdpIpcCommunicator interface {
	New() *UdpIpcConnection
	Listen(handler func(conn net.Conn)) error
}

// NewConnection creates a new UdpIpcConnection instance.
// socketPath - path to the socket file.
func NewConnection(udpAddr *net.UDPAddr) *UdpIpcConnection {
	return &UdpIpcConnection{Addr: udpAddr}
}

// NewConnectionFromIpPort creates a new UdpIpcConnection instance.
// socketPath - path to the socket file.
func NewConnectionFromIpPort(ip net.IP, port int) *UdpIpcConnection {
	return &UdpIpcConnection{Addr: &net.UDPAddr{
		IP:   ip,
		Port: port,
		Zone: "",
	}}
}

// Listen listens to the local socket.
// handler - a function that handles the connection.
func (ipcConnection *UdpIpcConnection) Listen(handler func(conn *net.UDPConn)) error {
	socket, err := net.ListenUDP("udp", ipcConnection.Addr)
	if err != nil {
		return err
	}

	for {
		conn, err := socket.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		go handler(conn)
	}
}
