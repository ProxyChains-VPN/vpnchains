package main

import (
	"log"
	"net"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/ipc_request/udp_ipc_request"
	"vpnchains/gopkg/vpn"
)

type packetOwner struct {
	pid int64
	fd  int32
}

var packets map[packetOwner][]byte = make(map[packetOwner][]byte)

func handleUdpIpcMessage(sockAddr *net.UDPAddr, requestPacket []byte, bufSize int, tunnel vpn.UdpTunnel) {
	requestType, err := ipc.GetRequestType(requestPacket)
	if err != nil {
		log.Println("ERROR PARSING", err)
		return
	}

	switch requestType {
	case "recvfrom":
		request, err := udp_ipc_request.RecvfromRequestFromBytes(requestPacket)
		if err != nil {
			log.Println("converting recvfrom request from bytearray", err)
			return
		}

		sa := udp_ipc_request.UnixIpPortToUDPAddr(uint32(request.SrcIp), uint16(request.SrcPort))
		log.Println("recvfrom sa", sa.IP, sa.Port)

		packet := packets[packetOwner{request.Pid, request.Fd}]
		var response udp_ipc_request.RecvfromResponse
		if packet == nil {
			log.Println("no packet for fd", request.Fd)
			response = udp_ipc_request.ErrorRecvfromResponse
		} else {
			response = udp_ipc_request.RecvfromResponse{
				BytesRead: int64(len(packet)),
				Msg:       packet,
			}
		}

		bytes, err := udp_ipc_request.RecvfromResponseToBytes(response)
		if err != nil {
			log.Println("error serializing response", err)
			return
		}

		udp, err := net.DialUDP("udp", nil, sockAddr)
		if err != nil {
			log.Println("error dialing local process", err)
			return
		}

		_, err = udp.Write(bytes)
		if err != nil {
			log.Println("error writing to local process", err)
			return
		}

	case "sendto":
		request, err := udp_ipc_request.SendtoRequestFromBytes(requestPacket)
		if err != nil {
			log.Println("error parsing request", err)
			return
		}

		sa := udp_ipc_request.UnixIpPortToUDPAddr(uint32(request.DestIp), uint16(request.DestPort))
		log.Println("sendto sa", sa.IP, sa.Port)

		conn, err := tunnel.Dial(sa)
		if err != nil {
			log.Println("error dialing", err)
			return
		}

		_, err = conn.Write(request.Msg[:request.MsgLen])
		if err != nil {
			return
		}

		buf := make([]byte, bufSize)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("error reading from conn", err)
			return
		}

		packets[packetOwner{request.Pid, request.Fd}] = buf[:n]
	default:
		log.Println("unknown request type:", requestType)
		return
	}
}
