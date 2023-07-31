package main

import (
	"log"
	"net"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/ipc_request/udp_ipc_request"
	"vpnchains/gopkg/vpn"
)

type PacketOwner struct {
	pid int64
	fd  int32
}

type Packet struct {
	bytes []byte
	ip    int32
	port  uint16
}

var packets *PacketsBuffer = NewPacketsBuffer()

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

		log.Println("recvfrom request", "pid/fd", request.Fd, request.Pid)

		var response udp_ipc_request.RecvfromResponse

		packet := packets.WaitForPacket(PacketOwner{request.Pid, request.Fd})

		if packet == nil {
			response = udp_ipc_request.RecvfromResponse{
				BytesRead: -1,
				Msg:       []byte{},
			}
		} else {
			response = udp_ipc_request.RecvfromResponse{
				BytesRead: int64(len(packet.bytes)),
				Msg:       packet.bytes,
				SrcIp:     packet.ip,
				SrcPort:   packet.port,
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
		log.Println("sendto sa", sa.IP, sa.Port, "pid/fd", request.Pid, request.Fd)

		conn, err := tunnel.Dial(sa)
		if err != nil {
			log.Println("error dialing", err)
			return
		}

		buf := make([]byte, bufSize)

		_, err = conn.Write(request.Msg[:request.MsgLen])
		if err != nil {
			return
		}

		go func() {
			for {
				n, err := conn.Read(buf)
				if err != nil {
					log.Println("error reading from conn", err)
					return
				}

				log.Println("read", n, "bytes")

				recvPacket := &Packet{
					bytes: buf[:n],
					ip:    request.DestIp,
					port:  request.DestPort,
				}

				packets.PushPacket(PacketOwner{request.Pid, request.Fd}, recvPacket)
			}
		}()

	default:
		log.Println("unknown request type:", requestType)
		return
	}
}
