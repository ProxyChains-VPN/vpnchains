package wireguard

import (
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"net"
	"net/netip"
)

type WireguardTunnel struct {
	Dev      *device.Device // TODO а оно пригодится???
	Net      *netstack.Net
	TcpFdMap map[int32]*net.Conn
	//config   *WireguardConfig
} // TODO инкапсулировать инкапсулируемое

func NewWireguardTunnel(localAddresses, dnsAddresses []netip.Addr, mtu int, uapiConfig string) (*WireguardTunnel, error) {
	tun, tnet, err := netstack.CreateNetTUN(localAddresses, dnsAddresses, mtu)
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

	tunnel := &WireguardTunnel{
		Dev:      dev,
		Net:      tnet,
		TcpFdMap: make(map[int32]*net.Conn),
	}

	return tunnel, nil
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

	uapi, err := config.Uapi()
	if err != nil {
		return nil, err
	}

	return NewWireguardTunnel(localAddresses, dnsAddresses, mtu, uapi)
}

func (tunnel *WireguardTunnel) CloseTunnel() {
	tunnel.Dev.Close()
}
