package wireguard

import (
	"encoding/base64"
	"encoding/hex"
	"gopkg.in/ini.v1"
	"net/netip"
	"strings"
)

// WireguardConfig represents a wireguard config file, that
// is in the format of the following example:
//
// [Interface]
// PrivateKey: ...
// Address: ...
// DNS: ...
//
// [Peer]
// Endpoint: ...
// AllowedIPs: ...
// PublicKey: ...
// PresharedKey: ...
//
// In fact, wireguard config files are actually INI files,
// so we use the gopkg.in/ini.v1 package to parse them.
//
// Todo: fix problems with lowercase and uppercase letters.
//
// See https://www.wireguard.com/quickstart/ for more information.
type WireguardConfig struct {
	Interface struct {
		PrivateKey string   `ini:"PrivateKey"`
		Address    []string `ini:"Address"`
		DNS        []string `ini:"DNS"`
	}
	Peer struct {
		Endpoint     string   `ini:"Endpoint"`
		AllowedIPs   []string `ini:"AllowedIPs"`
		PublicKey    string   `ini:"PublicKey"`
		PresharedKey string   `ini:"PresharedKey"`
	}
}

// WireguardConfigFromFile parses a wireguard config file.
// path - path to the config file.
func WireguardConfigFromFile(path string) (*WireguardConfig, error) {
	var config WireguardConfig

	err := ini.MapTo(&config, path)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func decodeKey(key string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(decodedKey), nil
}

func splitAddress(address string) (ip, port string, err error) {
	arr := strings.SplitN(address, "/", 3)
	if len(arr) != 2 {
		return "", "", nil
	}
	return arr[0], arr[1], nil
}

func (config *WireguardConfig) addressStringToNetipAddr() ([]netip.Addr, error) { // TODO rename
	var res []netip.Addr
	for _, addr := range config.Interface.Address {
		ip, _, err := splitAddress(addr)

		if err != nil {
			return nil, err
		}

		netipAddr, err := netip.ParseAddr(ip)
		if err != nil {
			return nil, err
		}
		res = append(res, netipAddr)
	}
	return res, nil
}

func (config *WireguardConfig) dnsStringToNetipAddr() ([]netip.Addr, error) { // TODO rename
	var res []netip.Addr
	for _, addr := range config.Interface.DNS {
		netipAddr, err := netip.ParseAddr(addr)
		if err != nil {
			return nil, err
		}
		res = append(res, netipAddr)
	}
	return res, nil
}

// UapiString returns a string that can be used while configurating the tunnel (IpcSet) todo.
func (config *WireguardConfig) UapiString() (string, error) {
	privateKeyDecoded, err := decodeKey(config.Interface.PrivateKey)
	if err != nil {
		return "", err
	}
	publicKeyDecoded, err := decodeKey(config.Peer.PublicKey)
	if err != nil {
		return "", err
	}
	presharedKeyDecoded, err := decodeKey(config.Peer.PresharedKey)
	if err != nil {
		return "", err
	}

	var res string
	res += `private_key=` + privateKeyDecoded + "\n"
	res += `public_key=` + publicKeyDecoded + "\n"
	if presharedKeyDecoded != "" {
		res += `preshared_key=` + presharedKeyDecoded + "\n"
	}

	for _, addr := range config.Peer.AllowedIPs {
		res += `allowed_ip=` + addr + "\n"
	}

	res += `endpoint=` + config.Peer.Endpoint + "\n"
	return res, nil
}
