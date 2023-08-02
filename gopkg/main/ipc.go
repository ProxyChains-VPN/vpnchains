package main

import (
	"log"
	"net"
	"vpnchains/gopkg/ipc/tcp_ipc"
	"vpnchains/gopkg/vpn"
)

func startIpcWithSubprocess(tcpTunnel vpn.TcpTunnel, port int, bufSize int) {
	tcpConn := tcp_ipc.NewConnectionFromIpPort(net.IPv4(127, 0, 0, 1), port)

	err := tcpConn.Listen(
		func(sockConn *net.TCPConn) {
			handleTcpIpcMessage(sockConn, bufSize, tcpTunnel)
		},
	)
	if err != nil {
		log.Fatalln("unable to start listening", err)
	}

	if err != nil {
		log.Fatalln("unable to start udp reading", err)
	}
}
