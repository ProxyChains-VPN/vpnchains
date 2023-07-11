package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/vpn/wireguard"
	"syscall"
)

var errorConnectResponse = ipc.ConnectResponse{ResultCode: -1}
var errorConnectResponseBytes, _ = bson.Marshal(errorConnectResponse)

func (handler *RequestHandler) processConnect(request *ipc.ConnectRequest) (*ipc.ConnectResponse, error) {
	//log.Println(request)

	/*var sa syscall.SockaddrInet4
	sa.Addr[0] = byte(request.Ip >> 24)
	sa.Addr[1] = byte(request.Ip >> 16)
	sa.Addr[2] = byte(request.Ip >> 8)
	sa.Addr[3] = byte(request.Ip)
	sa.Port = int(request.Port)
	connect4(request.SockFd, &sa)*/

	return &ipc.ConnectResponse{0}, nil
}
