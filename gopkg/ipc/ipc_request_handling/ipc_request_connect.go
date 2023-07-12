package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"syscall"
	"vpnchains/gopkg/ipc"
)

var errorConnectResponse = ipc.ConnectResponse{ResultCode: -1}
var errorConnectResponseBytes, _ = bson.Marshal(errorConnectResponse)

func unixIpToSockaddr(ip uint32, port uint16) syscall.SockaddrInet4 {
	var sa syscall.SockaddrInet4
	sa.Addr[0] = byte(ip >> 24)
	sa.Addr[1] = byte(ip >> 16)
	sa.Addr[2] = byte(ip >> 8)
	sa.Addr[3] = byte(ip)
	sa.Port = int(port)

	log.Println(sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3])
	return sa
}

func (handler *RequestHandler) processConnect(request *ipc.ConnectRequest) (*ipc.ConnectResponse, error) {
	log.Println(request)
	sa := unixIpToSockaddr(request.Ip, request.Port)
	err := handler.tunnel.Connect(request.SockFd, &sa)
	if err != nil {
		return &errorConnectResponse, err
	}

	return &ipc.ConnectResponse{0}, nil //success
}
