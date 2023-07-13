package ipc_request

import (
	"go.mongodb.org/mongo-driver/bson"
	"syscall"
)

var ErrorConnectResponse = ConnectResponse{ResultCode: -1}
var SuccConnectResponse = ConnectResponse{ResultCode: 0}

func IpPortToSockaddr(ip uint32, port uint16) syscall.SockaddrInet4 {
	var sa syscall.SockaddrInet4
	sa.Addr[3] = byte(ip >> 24)
	sa.Addr[2] = byte(ip >> 16)
	sa.Addr[1] = byte(ip >> 8)
	sa.Addr[0] = byte(ip)
	sa.Port = int(port)

	return sa
}

func (handler *RequestHandler) ConnectRequestFromBytes(requestBytes []byte) (*ConnectRequest, error) {
	var connectRequest ConnectRequest
	err := bson.Unmarshal(requestBytes, &connectRequest)
	if err != nil {
		return nil, err
	}
	return &connectRequest, nil
}

func (handler *RequestHandler) ConnectResponseToBytes(response ConnectResponse) ([]byte, error) {
	return bson.Marshal(response) // todo err handling
}
