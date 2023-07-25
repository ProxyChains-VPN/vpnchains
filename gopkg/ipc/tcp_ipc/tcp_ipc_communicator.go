package tcp_ipc

import (
	"log"
	"net"
)

// TcpIpcConnection A struct that represents an IPC connection.
// Addr - a net.TCPAddr instance.
type TcpIpcConnection struct {
	Addr *net.TCPAddr
}

// TcpIpcCommunicator An interface that represents an IPC communicator.
// Listen(handler func(conn *net.TCPConn)) - listens to the local socket.
type TcpIpcCommunicator interface {
	Listen(handler func(conn *net.TCPConn)) error
}

// NewConnection creates a new TcpIpcConnection instance.
// socketPath - path to the socket file.
func NewConnection(tcpAddr *net.TCPAddr) TcpIpcCommunicator {
	return &TcpIpcConnection{Addr: tcpAddr}
}

// NewConnectionFromIpPort creates a new TcpIpcConnection instance.
// socketPath - path to the socket file.
func NewConnectionFromIpPort(ip net.IP, port int) TcpIpcCommunicator {
	return &TcpIpcConnection{Addr: &net.TCPAddr{
		IP:   ip,
		Port: port,
		Zone: "",
	}}
}

// Listen listens to the local socket.
// handler - a function that handles the connection.
func (ipcConnection *TcpIpcConnection) Listen(handler func(conn *net.TCPConn)) error {
	socket, err := net.ListenTCP("tcp", ipcConnection.Addr)
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
