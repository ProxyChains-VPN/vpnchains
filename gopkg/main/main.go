package main

import (
	"flag"
	"log"
	"os"
	"strconv"
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
// If the path is neither specified in flags nor in the environment (InjectedLibPathEnvVar),
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

// DefaultWireguardConfigPath is the default path to the wireguard config.
// If the path is not specified in flags, this path (wg0.conf, I guess) will be used.
const DefaultWireguardConfigPath = "wg0.conf"

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return defaultValue
	}
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err != nil {
			return intValue
		}
	}
	return defaultValue
}

func main() {
	envMtu := getEnvOrDefaultInt(MtuEnvVar, DefaultMtu)
	envBufSize := getEnvOrDefaultInt(BufSizeEnvVar, DefaultBufSize)
	envIpcServerPort := getEnvOrDefaultInt(IpcServerPortEnvVar, DefaultIpcServerPort)
	envInjectedLibPath := getEnvOrDefault(InjectedLibPathEnvVar, DefaultInjectedLibPath)

	mtu := flag.Int("mtu", envMtu, "mtu for the wireguard tunnel")
	bufSize := flag.Int("buf", envBufSize, "size of the buffer used for reading from the socket")
	ipcServerPort := flag.Int("port", envIpcServerPort, "port for the IPC server [0, 65535]")
	injectedLibPath := flag.String("injected-lib-path", envInjectedLibPath, "path to the injected library")

	wireguardConfigPath := flag.String("config", DefaultWireguardConfigPath, "path to the wireguard config")

	flag.Parse()

	values := flag.Args()

	if len(values) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if *ipcServerPort < 0 || *ipcServerPort > 65535 {
		log.Fatalln("Invalid port number. Must be in range [0, 65535].")
	}
	commandPath := values[0]
	commandArgs := values[1:]

	err := os.Setenv(IpcServerPortEnvVar, strconv.Itoa(*ipcServerPort))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	config, err := wireguard.WireguardConfigFromFile(*wireguardConfigPath)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	tunnel, err := wireguard.TunnelFromConfig(config, *mtu)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer tunnel.CloseTunnel()

	cmd := ipc.CreateCommandWithInjectedLibrary(*injectedLibPath, commandPath, commandArgs)

	ready := make(chan struct{})
	go startIpcWithSubprocess(ready, tunnel, *ipcServerPort, *bufSize)

	<-ready
	err = cmd.Run()
	if err != nil {
		log.Fatalln("subprocess says,", err)
		os.Exit(1)
	}

	tunnel.CloseTunnel()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
