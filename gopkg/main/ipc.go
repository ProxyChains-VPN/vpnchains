package main

import (
	"log"
	"net"
	"vpnchains/gopkg/ipc/tcp_ipc"
	"vpnchains/gopkg/ipc/udp_ipc"
	"vpnchains/gopkg/vpn"
)

func startIpcWithSubprocess(tcpTunnel vpn.TcpTunnel, udpTunnel vpn.UdpTunnel, port int, bufSize int) {
	tcpConn := tcp_ipc.NewConnectionFromIpPort(net.IPv4(127, 0, 0, 1), port)
	udpConn := udp_ipc.NewConnectionFromIpPort(net.IPv4(127, 0, 0, 1), port, bufSize)

	err := tcpConn.Listen(
		func(sockConn *net.TCPConn) {
			handleTcpIpcMessage(sockConn, bufSize, tcpTunnel)
		},
	)
	if err != nil {
		log.Fatalln("unable to start listening", err)
	}

	log.Println("start udp reading")
	err = udpConn.ReadLoop(
		func(sockAddr *net.UDPAddr, requestPacket []byte) {
			handleUdpIpcMessage(sockAddr, requestPacket, bufSize, udpTunnel)
		},
	)

	if err != nil {
		log.Fatalln("unable to start udp reading", err)
	}
}
