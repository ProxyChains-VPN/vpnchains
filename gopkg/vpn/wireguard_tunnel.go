package vpn

import (
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"net"
	"net/netip"
	"syscall"
)

type WireguardTunnel struct {
	dev      *device.Device // TODO а оно пригодится???
	net      *netstack.Net
	tcpFdMap map[int]*net.Conn
	//config   *WireguardConfig
} // TODO инкапсулировать инкапсулируемое

func NewWireguardTunnel(localAddresses, dnsAddresses []netip.Addr, mtu int, uapiConfig string) (*WireguardTunnel, error) {
	tun, tnet, err := netstack.CreateNetTUN(dnsAddresses, localAddresses, 1420)
	if err != nil {
		return nil, err
	}
	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(device.LogLevelVerbose, ""))

	err = dev.IpcSet(uapiConfig)
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
	}, nil
}

func WireguardTunnelFromConfig(config *WireguardConfig, mtu int) (*WireguardTunnel, error) {
	localAddresses, err := config.AddressStringToNetipAddr()
	if err != nil {
		return nil, err
	}

	dnsAddresses, err := config.DnsStringToNetipAddr()
	if err != nil {
		return nil, err
	}

	uapi, err := config.UapiConfig()
	if err != nil {
		return nil, err
	}

	return NewWireguardTunnel(localAddresses, dnsAddresses, mtu, uapi)
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
