package vpn

import (
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"net"
	"syscall"
)

var Mtu = 1420

type WireguardTunnel struct {
	dev      *device.Device // TODO а оно пригодится
	net      *netstack.Net
	tcpFdMap map[int]*net.Conn
	config   *WireguardConfig
}

func StartWireguardTunnel(config *WireguardConfig) (*WireguardTunnel, error) {
	localAddresses, err := config.AddressStringToNetipAddr()
	if err != nil {
		return nil, err
	}

	dnsAddresses, err := config.DnsStringToNetipAddr()
	if err != nil {
		return nil, err
	}

	tun, tnet, err := netstack.CreateNetTUN(dnsAddresses, localAddresses, 1420)
	if err != nil {
		return nil, err
	}
	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(device.LogLevelVerbose, ""))

	uapi, err := config.UapiConfig()
	if err != nil {
		return nil, err
	}

	err = dev.IpcSet(uapi)
	if err != nil {
		return nil, err
	}
	err = dev.Up()
	if err != nil {
		return nil, err
	}

	return &WireguardTunnel{
		dev:      dev,
		net:      tnet,
		tcpFdMap: make(map[int]*net.Conn),
		config:   config,
	}, nil
}

func (tunnel *WireguardTunnel) Close() {
	tunnel.dev.Close()
}

func (tunnel *WireguardTunnel) Connect(fd int, sa syscall.Sockaddr) (err error) {
	return nil
}

func (tunnel *WireguardTunnel) Read(fd int, p []byte) (n int, err error) {
	return 0, nil
}

func (tunnel *WireguardTunnel) Write(fd int, p []byte) (n int, err error) {
	return 0, nil
}

func (tunnel *WireguardTunnel) GetFd() int {
	return 0
}
