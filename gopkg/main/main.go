package main

import (
	"fmt"
	"log"
	"os"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/vpn/wireguard"
)

const DefaultSockAddr = "/tmp/vpnchains.socket"
const InjectedLibPath = "/usr/lib/libvpnchains_inject.so"
const BufSize = 100500
const Mtu = 1420

func errorMsg(path string) string {
	return "Usage: " + path + " <config> " +
		"<command> [command args...]"
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println(errorMsg(args[0]))
		os.Exit(0)
	}

	wireguardConfigPath := args[1]
	commandPath := args[2]
	commandArgs := args[3:]

	config, err := wireguard.WireguardConfigFromFile(wireguardConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	tunnel, err := wireguard.TunnelFromConfig(config, Mtu)
	if err != nil {
		log.Fatalln(err)
	}
	defer tunnel.CloseTunnel()

	cmd := ipc.CreateCommandWithInjectedLibrary(InjectedLibPath, commandPath, commandArgs)

	ready := make(chan struct{})
	go startIpcWithSubprocess(ready, tunnel)

	<-ready
	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalln("subprocess says,", err)
	}

	tunnel.CloseTunnel()

	err = os.Remove(DefaultSockAddr)
	if err != nil {
		log.Println(err)
	}
}
