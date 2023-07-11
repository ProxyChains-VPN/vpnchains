package vpn

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
)

type connConfig struct {
	TunAddr    string `json:"tunAddr"`
	DnsAddr    string `json:"dnsAddr"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	ServerAddr string `json:"serverAddr"`
	ServerPort string `json:"serverPort"`
	AllowedIp  string `json:"allowedIp"`
	Network    string `json:"network"` //TODO: make it useful:)
}

func newConnConfig(file *os.File) (c *connConfig, err error) {
	buff := make([]byte, 4096, 4096)
	n, err := file.Read(buff)
	if err != nil {
		return nil, err
	}
	c = new(connConfig)
	err = json.Unmarshal(buff[:n], c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var (
	absPath, _ = filepath.Abs("./overrides/config.json")
	file, _    = os.Open(absPath)
	config, _  = newConnConfig(file)
	conns      = make(map[int]*net.Conn)
)
