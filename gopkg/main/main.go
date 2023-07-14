package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/ipc/ipc_request"
	"vpnchains/gopkg/vpn"
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

func handleIpc(ready chan struct{}, tunnel vpn.Tunnel) {
	err := os.Remove(DefaultSockAddr)
	if err != nil {
		log.Println(err)
	}

	var buf = make([]byte, BufSize)

	conn := ipc.NewConnection(DefaultSockAddr)
	requestHandler := ipc_request.NewRequestHandler(tunnel) // todo rename???

	ipcConnectionHandler := func(sockConn net.Conn) {
		n, err := sockConn.Read(buf)
		requestBuf := buf[:n]

		if err != nil {
			log.Fatalln(err)
		}

		requestType, err := requestHandler.ParseRequestType(requestBuf)

		switch requestType {
		case "connect":
			log.Println("connect")
			request, err := requestHandler.ConnectRequestFromBytes(requestBuf)

			sa := ipc_request.IpPortToSockaddr(uint32(request.Ip), request.Port) // todo сделать норм архитектуру и доделать метод
			// todo по хорошему надо хотяб отдельно функцию сделать аля handleConnect
			endpointConn, err := tunnel.Connect(request.SockFd, &sa)
			if err != nil {
				bytes, _ := requestHandler.ConnectResponseToBytes(ipc_request.ErrorConnectResponse)
				sockConn.Write(bytes)
			}

			bytes, _ := requestHandler.ConnectResponseToBytes(ipc_request.SuccConnectResponse)
			sockConn.Write(bytes)

			go func() {
				buf := make([]byte, BufSize)
				for {
					n, err := sockConn.Read(buf)
					if err != nil {
						log.Println("read from client", err)
						continue
					}
					log.Println(string(buf[:n]))
					_, err = endpointConn.Write(buf[:n])
					if err != nil {
						log.Println("write to server", err)
						continue
					}
				}
			}()

			go func() {
				buf := make([]byte, BufSize)
				for {
					n, err := endpointConn.Read(buf)
					if err != nil {
						log.Println("read from server", err)
						continue
					}
					_, err = sockConn.Write(buf[:n])
					if err != nil {
						log.Println("write to client", err)
						continue
					}
				}
			}()
			log.Println("connect ended")
		default:
			log.Println("Unknown request type:", requestType)
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

	tunnel, err := wireguard.TunnelFromConfig(config, Mtu)
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
		log.Fatalln("subprocess says,", err)
	}

	tunnel.CloseTunnel()

	err = os.Remove(DefaultSockAddr)
	if err != nil {
		log.Println(err)
	}
}
