package ipc

import (
	"log"
	"net"
)

type IpcConnection struct {
	SocketPath string
}

type IpcCommunicator interface {
	New() *IpcConnection
	Listen() error
}

func NewConnection(socketPath string) *IpcConnection {
	return &IpcConnection{SocketPath: socketPath}
}

func (ipcConnection *IpcConnection) Listen(handler func(conn net.Conn)) error {
	socket, err := net.Listen("unix", ipcConnection.SocketPath)
	if err != nil {
		log.Println("ladskf")
		return err
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			return err
		}

		go func() {
			handler(conn)
			conn.Close()
		}()
	}
}
