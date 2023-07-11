package wireguard

import (
	"encoding/base64"
	"encoding/hex"
	"gopkg.in/ini.v1"
	"net/netip"
)

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

func (config *WireguardConfig) AddressStringToNetipAddr() ([]netip.Addr, error) { // TODO rename
	var res []netip.Addr
	for _, addr := range config.Interface.Address {
		netipAddr, err := netip.ParseAddr(addr)
		if err != nil {
			return nil, err
		}
		res = append(res, netipAddr)
	}
	return res, nil
}

func (config *WireguardConfig) DnsStringToNetipAddr() ([]netip.Addr, error) { // TODO rename
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

func (config *WireguardConfig) UapiConfig() (string, error) {
	privateKeyDecoded, err := decodeKey(config.Interface.PrivateKey)
	if err != nil {
		return "", err
	}
	publicKeyDecoded, err := decodeKey(config.Peer.PublicKey)
	if err != nil {
		return "", err
	}

	var res string
	res += `private_key=` + privateKeyDecoded + "\n"
	res += `public_key=` + publicKeyDecoded + "\n"

	for _, addr := range config.Peer.AllowedIPs {
		res += `allowed_ip=` + addr + "\n"
	}

	res += `endpoint=` + config.Peer.Endpoint + "\n"
	return res, nil
}
