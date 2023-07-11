package overrides

import (
	"bytes"
	"fmt"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"net/netip"
	"strconv"
	"syscall"
	"vpnchains/gopkg/vpn"
)

func Connect(fd int, sa syscall.Sockaddr) (err error) {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return connect4(fd, sa)
	case *syscall.SockaddrInet6:
		return connect6(fd, sa)
	case *syscall.SockaddrUnix:
		return syscall.Connect(fd, sa)
	}
	return nil
}

func connect4(fd int, sa *syscall.SockaddrInet4) (err error) {
	if bytes.Equal([]byte{sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3]}, []byte{127, 0, 0, 1}) {
		return syscall.Connect(fd, sa)
	}
	tun, tnet, err := netstack.CreateNetTUN(
		[]netip.Addr{netip.MustParseAddr(config.TunAddr)},
		[]netip.Addr{netip.MustParseAddr(config.DnsAddr)},
		1420)
	if err != nil {
		return err
	}
	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(device.LogLevelVerbose, ""))
	err = dev.IpcSet(fmt.Sprintf(`private_key=%s
public_key=%s
endpoint=%s
allowed_ip=%s
`, config.PrivateKey, config.PublicKey, config.ServerAddr+":"+config.ServerPort, config.AllowedIp))
	if err != nil {
		return err
	}
	err = dev.Up()
	if err != nil {
		return err
	}
	address := strconv.Itoa(int(sa.Addr[0])) + "." + strconv.Itoa(int(sa.Addr[1])) + "." + strconv.Itoa(int(sa.Addr[2])) + "." + strconv.Itoa(int(sa.Addr[3])) + ":" + strconv.Itoa(sa.Port)
	socket, err := tnet.Dial(config.Network, address)
	if err != nil {
		return err
	}
	vpn.tcpConnsMap[fd] = &socket
	return nil
}

func connect6(fd int, sa *syscall.SockaddrInet6) (err error) {
	return nil
}
