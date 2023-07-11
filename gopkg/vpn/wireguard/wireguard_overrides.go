package wireguard

import "syscall"

func (tunnel *WireguardTunnel) Connect(fd int, sa syscall.Sockaddr) (err error) {
	return nil
}

func (tunnel *WireguardTunnel) Read(fd int, p []byte) (n int, err error) {
	return 0, nil
}

func (tunnel *WireguardTunnel) Write(fd int, p []byte) (n int, err error) {
	return 0, nil
}
