package main

import (
	"flag"
	"log"
	"os"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/vpn/wireguard"
)

// DefaultIpcServerPort is the default port for the IPC server.
// If the port is neither specified in flags nor in the environment (IpcServerPortEnvVar),
// this port (45454, I guess) will be used.
const DefaultIpcServerPort = 45454

// IpcServerPortEnvVar is the name of the environment variable that
// contains the port for the IPC server.
const IpcServerPortEnvVar = "VPNCHAINS_IPC_SERVER_PORT"

// DefaultInjectedLibPath is the default path to the injected library.
// If the path is not specified in the environment (InjectedLibPathEnvVar),
// this path (/usr/lib/libvpnchains_inject.so, I guess) will be used.
const DefaultInjectedLibPath = "/usr/lib/libvpnchains_inject.so"

// InjectedLibPathEnvVar is the name of the environment variable that
// contains the path to the injected library.
const InjectedLibPathEnvVar = "VPNCHAINS_INJECT_LIB_PATH"

// DefaultBufSize is the default size of the buffer used for reading from the socket.
// If the size is neither specified in flags nor in the environment (BufSizeEnvVar),
// this size (100500, I guess) will be used.
const DefaultBufSize = 100500

// BufSizeEnvVar is the name of the environment variable that
// contains the size of the buffer used for reading from the socket.
const BufSizeEnvVar = "VPNCHAINS_BUF_SIZE"

// DefaultMtu is the default mtu for the wireguard tunnel.
// If the mtu is neither specified in flags nor in the environment (MtuEnvVar),
// this amount (1420, I guess) will be used.
const DefaultMtu = 1420

// MtuEnvVar is the name of the environment variable that
// contains the mtu for the wireguard tunnel.
const MtuEnvVar = "VPNCHAINS_MTU"

//func errorMsg(path string) string {
//	return "Usage: " + path + " <config> " +
//		"<command> [command args...]"
//}

func parseServer() {

}

func main() {
	//args := os.Args
	//if len(args) < 3 {
	//	fmt.Println(errorMsg(args[0]))
	//	os.Exit(0)
	//}

	mtu := flag.Int("--mtu", -1, "mtu for the wireguard tunnel")
	injectedLibPath := flag.String("--inject-lib", "", "path to the injected library")

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
