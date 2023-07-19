package ipc_request

import (
	"go.mongodb.org/mongo-driver/bson"
	"net"
)

var ErrorConnectResponse = ConnectResponse{ResultCode: -1}

var SuccConnectResponse = ConnectResponse{ResultCode: 0}

func UnixIpPortToTCPAddr(unixIp uint32, port uint16) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   net.IPv4(byte(unixIp), byte(unixIp>>8), byte(unixIp>>16), byte(unixIp>>24)),
		Port: int(port),
		Zone: "",
	}
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
	return bson.Marshal(response)
}
