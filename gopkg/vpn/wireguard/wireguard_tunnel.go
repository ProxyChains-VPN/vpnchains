package wireguard

import (
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"net/netip"
)

// WireguardTunnel An interface that represents a VPN tunnel.
type WireguardTunnel struct {
	dev *device.Device // TODO а оно пригодится???
	net *netstack.Net
} // TODO инкапсулировать инкапсулируемое

// NewTunnel creates a new tunnel.
// localAddresses - local addresses of the tunnel.
// dnsAddresses - dns addresses of the tunnel.
// mtu - mtu of the tunnel.
// uapiConfig - uapi config of the tunnel.
func NewTunnel(localAddresses, dnsAddresses []netip.Addr, mtu int, uapiConfig string) (*WireguardTunnel, error) {
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
		dev: dev,
		net: tnet,
	}

	return tunnel, nil
}

// TunnelFromConfig creates a tunnel from a wireguard config.
// config - wireguard config.
// mtu - mtu of the tunnel.
func TunnelFromConfig(config *WireguardConfig, mtu int) (*WireguardTunnel, error) {
	localAddresses, err := config.addressStringToNetipAddr()
	if err != nil {
		return nil, err
	}

	dnsAddresses, err := config.dnsStringToNetipAddr()
	if err != nil {
		return nil, err
	}

	uapi, err := config.UapiString()
	if err != nil {
		return nil, err
	}

	return NewTunnel(localAddresses, dnsAddresses, mtu, uapi)
}

// CloseTunnel closes the tunnel.
func (tunnel *WireguardTunnel) CloseTunnel() {
	tunnel.dev.Close()
}
