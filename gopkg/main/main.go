package main

import (
	"fmt"
	"log"
	"os"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/vpn/wireguard"
)

// DefaultIpcServerPort is the default port for the IPC server.
// If the port is neither specified in flags nor in the environment (VPNCHAINS_IPC_SERVER_PORT),
// this port (45454, I guess) will be used.
const DefaultIpcServerPort = 45454

// DefaultInjectedLibPath is the default path to the injected library.
// If the path is not specified in the environment (VPNCHAINS_INJECT_LIB_PATH),
// this path (/usr/lib/libvpnchains_inject.so, I guess) will be used.
const DefaultInjectedLibPath = "/usr/lib/libvpnchains_inject.so"

// DefaultBufSize is the default size of the buffer used for reading from the socket.
// If the size is neither specified in flags nor in the environment (VPNCHAINS_BUF_SIZE),
// this size (100500, I guess) will be used.
const DefaultBufSize = 100500

// DefaultMtu is the default mtu for the wireguard tunnel.
// If the mtu is neither specified in flags nor in the environment (VPNCHAINS_MTU),
// this amount (1420, I guess) will be used.
const DefaultMtu = 1420

func errorMsg(path string) string {
	return "Usage: " + path + " <config> " +
		"<command> [command args...]"
}

func getIpcServerPort() int {

}

func environParseServerPort() {

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
