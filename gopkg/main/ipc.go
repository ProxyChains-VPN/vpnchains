package main

import (
	"log"
	"net"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/ipc/tcp_ipc"
	"vpnchains/gopkg/ipc/udp_ipc"
	"vpnchains/gopkg/ipc_request/tcp_ipc_request"
	"vpnchains/gopkg/vpn"
)

func handleUdpIpcMessage(sockConn *net.UDPConn, requestPacket []byte, bufSize int, tunnel vpn.UdpTunnel) {
	requestType, err := ipc.GetRequestType(requestPacket)
	if err != nil {
		log.Println("ERROR PARSING", err)
		return
	}

	switch requestType {
	case "recvfrom":
	}
}

func handleTcpIpcMessage(sockConn *net.TCPConn, bufSize int, tunnel vpn.TcpTunnel) {
	buf := make([]byte, bufSize)
	n, err := sockConn.Read(buf)
	requestBuf := buf[:n]

	if err != nil {
		log.Fatalln(err)
	}

	requestType, err := ipc.GetRequestType(requestBuf)
	if err != nil {
		log.Println("ERROR PARSING", err)
		return
	}

	switch requestType {
	case "connect":
		request, err := tcp_ipc_request.ConnectRequestFromBytes(requestBuf)
		if err != nil {
			log.Println("eRROR PARSING", err)
			return
		}

		sa := tcp_ipc_request.UnixIpPortToTCPAddr(uint32(request.Ip), request.Port)
		log.Println("connect to sa", sa.IP, sa.Port)
		endpointConn, err := tunnel.Connect(sa)
		if err != nil {
			log.Println("ERROR CONNECTING", err)
			bytes, _ := tcp_ipc_request.ConnectResponseToBytes(tcp_ipc_request.ErrorConnectResponse)
			sockConn.Write(bytes)
			return
		}

		// client writes to server
		go func() {
			buf := make([]byte, bufSize)
			for {
				n, err := sockConn.Read(buf)
				if err != nil {
					log.Println("read from client", err)
					log.Println("closing endpoint write and socket read")
					endpointConn.CloseWrite()
					sockConn.CloseRead()
					return
				}
				_, err = endpointConn.Write(buf[:n])
				if err != nil {
					log.Println("write to server", err)
					log.Println("closing endpoint write and socket read")
					endpointConn.CloseWrite()
					sockConn.CloseRead()
					return
				}
			}
		}()

		// server writes to client
		go func() {
			buf := make([]byte, bufSize)
			for {
				n, err := endpointConn.Read(buf)
				if err != nil {
					//if errors.Is(err, io.EOF) {
					log.Println("read from server", err)
					log.Println("closing endpoint read and socket write")
					endpointConn.CloseRead()
					sockConn.CloseWrite()
					return
				}
				//log.Println("READ FROM SERVER", string(buf[:n]))
				_, err = sockConn.Write(buf[:n]) // todo если что в несколько раз отправить?????
				if err != nil {
					log.Println("write to client", err)
					log.Println("closing endpoint read and socket write")
					endpointConn.CloseRead()
					sockConn.CloseWrite()
					return
				}
			}
		}()

		bytes, _ := tcp_ipc_request.ConnectResponseToBytes(tcp_ipc_request.SuccConnectResponse)
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

func startIpcWithSubprocess(ready chan struct{}, tunnel vpn.TcpTunnel, port int, bufSize int) {
	tcpConn := tcp_ipc.NewConnectionFromIpPort(net.IPv4(127, 0, 0, 1), port)
	udpConn := udp_ipc.NewConnectionFromIpPort(net.IPv4(127, 0, 0, 1), port, bufSize)

	ready <- struct{}{}
	err := tcpConn.Listen(
		func(sockConn *net.TCPConn) {
			handleTcpIpcMessage(sockConn, bufSize, tunnel)
		},
	)
	if err != nil {
		log.Println("unable to start listening", err)
		log.Fatalln(err)
	}

	err = udpConn.ReadLoop()
}
