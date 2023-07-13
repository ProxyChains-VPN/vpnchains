package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"syscall"
	"vpnchains/gopkg/so_ipc"
)

var errorConnectResponse = so_ipc.ConnectResponse{ResultCode: -1}
var errorConnectResponseBytes, _ = bson.Marshal(errorConnectResponse)

func ipPortToSockaddr(ip uint32, port uint16) syscall.SockaddrInet4 {
	var sa syscall.SockaddrInet4
	sa.Addr[3] = byte(ip >> 24)
	sa.Addr[2] = byte(ip >> 16)
	sa.Addr[1] = byte(ip >> 8)
	sa.Addr[0] = byte(ip)
	sa.Port = int(port)

	return sa
}

func (handler *RequestHandler) processConnect(request *so_ipc.ConnectRequest) (*so_ipc.ConnectResponse, error) {
	log.Println("connect", request)
	sa := ipPortToSockaddr(uint32(request.Ip), request.Port)
	err := handler.tunnel.Connect(request.SockFd, &sa)
	log.Println("connect ended")
	if err != nil {
		return &errorConnectResponse, err
	}

	return &so_ipc.ConnectResponse{0}, nil //success
}
