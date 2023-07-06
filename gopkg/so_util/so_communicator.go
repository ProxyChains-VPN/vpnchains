package so_util

import (
	"net"
	"os"
)

type IpcConnection struct {
	ProcessConn net.Conn
}

type IpcCommunicator interface {
	EstablishIpc(socketPath string) (*IpcConnection, error)
}

// на данный момент функция устанавливает связь однократно
func EstablishIpc(socketPath string) (*IpcConnection, error) {
	listener, err := net.ListenUnix(
		"unix",
		&net.UnixAddr{
			Name: socketPath,
			Net:  "unix",
		},
	)
	if err != nil {
		return nil, err
	}

	defer listener.Close()
	defer os.Remove(socketPath) // сделать многа сокетов?

	processConn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	return &IpcConnection{processConn}, nil
}
