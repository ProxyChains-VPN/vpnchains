package openVPN

import (
	"github.com/ooni/minivpn/vpn"
	"net"
)

type OvpnTunnel struct {
	dialer   vpn.TunDialer
	TcpFdMap map[int32]*net.Conn
}

func NewOvpnTunnel(opts *vpn.Options) (tun *OvpnTunnel, err error) {
	dialer := vpn.NewTunDialerFromOptions(opts)
	tunnel := &OvpnTunnel{
		dialer:   dialer,
		TcpFdMap: make(map[int32]*net.Conn),
	}

	return tunnel, nil
}

func OvpnTunnelFromConfig(filePath string) (tun *OvpnTunnel, err error) {
	opts, err := vpn.ParseConfigFile(filePath)
	if err != nil {
		return nil, err
	}

	return NewOvpnTunnel(opts)
}
