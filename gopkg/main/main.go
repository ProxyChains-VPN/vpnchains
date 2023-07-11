package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/ipc/ipc_request_handling"
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

func handleIpc(ready chan struct{}, tunnel *wireguard.WireguardTunnel) {
	err := os.Remove(DefaultSockAddr)
	if err != nil {
		log.Println(err)
	}

	var buf = make([]byte, BufSize)

	conn := ipc.NewConnection(DefaultSockAddr)
	requestHandler := ipc_request_handling.NewRequestHandler(tunnel) // todo rename???

	ipcConnectionHandler := func(conn net.Conn) {
		n, err := conn.Read(buf)
		requestBuf := buf[:n]

		if err != nil {
			log.Fatalln(err)
		}

		responseBuf, err := requestHandler.HandleRequest(requestBuf)
		if responseBuf == nil && err != nil {
			log.Fatalln(err) // вроде как невозможно
		} else if err != nil {
			log.Println(err, ". Returning error response.")
		}

		_, err = conn.Write(responseBuf)
		if err != nil {
			log.Println(err)
		}
	}

	ready <- struct{}{}
	err = conn.Listen(ipcConnectionHandler)
	if err != nil {
		log.Println("sldfadsf")
		log.Fatalln(err)
	}
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

	tunnel, err := wireguard.WireguardTunnelFromConfig(config, Mtu)
	if err != nil {
		log.Fatalln(err)
	}
	defer tunnel.CloseTunnel()

	cmd := ipc.CreateCommandWithInjectedLibrary(InjectedLibPath, commandPath, commandArgs)

	ready := make(chan struct{})
	go handleIpc(ready, tunnel)

	<-ready
	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalln("subprocess says, ", err)
	}

	tunnel.CloseTunnel()

	err = os.Remove(DefaultSockAddr)
	if err != nil {
		log.Println(err)
	}
}
