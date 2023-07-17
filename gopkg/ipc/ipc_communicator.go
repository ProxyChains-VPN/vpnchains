package ipc

import (
	"log"
	"net"
)

// IpcConnection A struct that represents an IPC connection.
// SocketPath - path to the socket file.
type IpcConnection struct {
	Addr *net.UnixAddr
}

// IpcCommunicator An interface that represents an IPC communicator.
// New() - creates a new IpcConnection instance.
// Listen(handler func(conn net.Conn)) - listens to the socket.
type IpcCommunicator interface {
	New() *IpcConnection
	Listen(handler func(conn net.Conn)) error
}

// NewConnection creates a new IpcConnection instance.
// socketPath - path to the socket file.
func NewConnection(socketPath string) *IpcConnection {
	return &IpcConnection{Addr: &net.UnixAddr{Name: socketPath}}
}

// Listen listens to the socket.
// handler - a function that handles the connection.
func (ipcConnection *IpcConnection) Listen(handler func(conn *net.UnixConn)) error {
	socket, err := net.ListenUnix("unix", ipcConnection.Addr)
	if err != nil {
		return err
	}

	for {
		conn, err := socket.AcceptUnix()
		if err != nil {
			log.Println(err)
			continue
		}

		go handler(conn)
	}
}
