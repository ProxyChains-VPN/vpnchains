package tcp_ipc

import (
	"log"
	"net"
)

// TcpIpcConnection A struct that represents an IPC connection.
// addr - a net.TCPAddr instance.
type TcpIpcConnection struct {
	addr *net.TCPAddr
}

// TcpIpcCommunicator An interface that represents an IPC communicator.
// Listen(handler func(conn *net.TCPConn)) - listens to the local socket.
type TcpIpcCommunicator interface {
	Listen(handler func(conn *net.TCPConn)) error
}

// NewConnection creates a new TcpIpcConnection instance.
// tcpAddr - a net.TCPAddr instance.
func NewConnection(tcpAddr *net.TCPAddr) TcpIpcCommunicator {
	return &TcpIpcConnection{addr: tcpAddr}
}

// NewConnectionFromIpPort creates a new TcpIpcConnection instance.
// ip - ip address.
// port - port.
func NewConnectionFromIpPort(ip net.IP, port int) TcpIpcCommunicator {
	return &TcpIpcConnection{addr: &net.TCPAddr{
		IP:   ip,
		Port: port,
		Zone: "",
	}}
}

// Listen listens to the local socket.
// handler - a function that handles the connection.
func (ipcConnection *TcpIpcConnection) Listen(handler func(conn *net.TCPConn)) error {
	socket, err := net.ListenTCP("tcp", ipcConnection.addr)
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
