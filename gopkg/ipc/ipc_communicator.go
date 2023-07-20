package ipc

import (
	"log"
	"net"
)

// IpcConnection A struct that represents an IPC connection.
// Addr - a net.TCPAddr instance.
type IpcConnection struct {
	Addr *net.TCPAddr
}

// IpcCommunicator An interface that represents an IPC communicator.
// New() - creates a new IpcConnection instance.
// Listen(handler func(conn net.Conn)) - listens to the local socket.
type IpcCommunicator interface {
	New() *IpcConnection
	Listen(handler func(conn net.Conn)) error
}

// NewConnection creates a new IpcConnection instance.
// socketPath - path to the socket file.
func NewConnection(tcpAddr *net.TCPAddr) *IpcConnection {
	return &IpcConnection{Addr: tcpAddr}
}

// NewConnectionFromIpPort creates a new IpcConnection instance.
// socketPath - path to the socket file.
func NewConnectionFromIpPort(ip net.IP, port int) *IpcConnection {
	return &IpcConnection{Addr: &net.TCPAddr{
		IP:   ip,
		Port: port,
		Zone: "",
	}}
}

// Listen listens to the local socket.
// handler - a function that handles the connection.
func (ipcConnection *IpcConnection) Listen(handler func(conn *net.TCPConn)) error {
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
