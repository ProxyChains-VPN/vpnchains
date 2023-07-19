package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/ipc/ipc_request"
	"vpnchains/gopkg/vpn"
)

func handleIpcMessage(sockConn *net.TCPConn, requestHandler *ipc_request.RequestHandler, buf []byte, tunnel vpn.Tunnel) {
	n, err := sockConn.Read(buf)
	requestBuf := buf[:n]

	if err != nil {
		log.Fatalln(err)
	}

	requestType, err := requestHandler.GetRequestType(requestBuf)

	switch requestType {
	case "connect":
		request, err := requestHandler.ConnectRequestFromBytes(requestBuf)
		if err != nil {
			log.Println("eRROR PARSING", err)
			return
		}

		sa := ipc_request.UnixIpPortToTCPAddr(uint32(request.Ip), request.Port)
		log.Println("connect to sa", sa.IP, sa.Port)
		endpointConn, err := tunnel.Connect(sa)
		if err != nil {
			log.Println("ERROR CONNECTING", err)
			bytes, _ := requestHandler.ConnectResponseToBytes(ipc_request.ErrorConnectResponse)
			sockConn.Write(bytes)
			return
		}

		// client writes to server
		go func() {
			buf := make([]byte, BufSize)
			for {
				n, err := sockConn.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						log.Println("closing endpoint write and socket read")
						sockConn.CloseRead()
						//endpointConn.CloseWrite()
						return
					} else {
						log.Println("read from client", err)
						continue
					}
				}
				//log.Println("READ FROM CLIENT", string(buf[:n]))
				_, err = endpointConn.Write(buf[:n])
				if err != nil {
					log.Println("write to server", err)
					continue
				}
			}
		}()

		// server writes to client
		go func() {
			buf := make([]byte, BufSize)
			for {
				n, err := endpointConn.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						log.Println("closing endpoint read and socket write")
						//endpointConn.CloseRead()
						sockConn.CloseWrite()
						return
					} else {
						log.Println("read from server", err)
						continue
					}
				}
				//log.Println("READ FROM SERVER", string(buf[:n]))
				_, err = sockConn.Write(buf[:n]) // todo если что в несколько раз отправить?????
				if err != nil {
					log.Println("write to client", err)
					continue
				}
			}
		}()

		bytes, _ := requestHandler.ConnectResponseToBytes(ipc_request.SuccConnectResponse)
		n, err = sockConn.Write(bytes)
		if err != nil {
			log.Println(err)
		}

		log.Println("connect ended")
	default:
		log.Println("Unknown request type:", requestType)
		return
	}
}

func startIpcWithSubprocess(ready chan struct{}, tunnel vpn.Tunnel) {
	err := os.Remove(DefaultSockAddr)
	if err != nil {
		log.Println(err)
	}

	var buf = make([]byte, BufSize)

	conn := ipc.NewConnection(DefaultSockAddr)
	requestHandler := ipc_request.NewRequestHandler(tunnel) // todo rename???

	ready <- struct{}{}
	err = conn.Listen(
		func(sockConn *net.TCPConn) {
			handleIpcMessage(sockConn, requestHandler, buf, tunnel)
		},
	)
	if err != nil {
		log.Println("unable to start listening", err)
		log.Fatalln(err)
	}
}
