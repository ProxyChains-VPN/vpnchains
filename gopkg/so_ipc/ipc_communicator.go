package so_ipc

import (
	"net"
)

// IpcConnection A struct that represents an IPC connection.
// SocketPath - path to the socket file.
type IpcConnection struct {
	SocketPath string
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
	return &IpcConnection{SocketPath: socketPath}
}

// Listen listens to the socket.
// handler - a function that handles the connection.
func (ipcConnection *IpcConnection) Listen(handler func(conn net.Conn)) error {
	socket, err := net.Listen("unix", ipcConnection.SocketPath)
	if err != nil {
		return err
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			return err
		}

		go func() {
			handler(conn)
			//conn.Close()
		}()
	}
}
