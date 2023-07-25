package udp_ipc_request

import (
	"go.mongodb.org/mongo-driver/bson"
	"net"
)

// ErrorConnectResponse A ConnectResponse instance that represents an unsuccessful connect syscall.
// Later errno also is to be added. (todo)
var ErrorConnectResponse = ConnectResponse{ResultCode: -1}

// SuccConnectResponse A ConnectResponse instance that represents a successful connect syscall.
// Later errno also is to be added. (todo)
var SuccConnectResponse = ConnectResponse{ResultCode: 0}

// UnixIpPortToTCPAddr A function that converts POSIX-style IP address (that is an unsigned 64-bit value)
// and port (that is an unsigned 16-bit value) to a net.TCPAddr instance.
// unixIp - POSIX-style IP address.
// port - port.
func UnixIpPortToTCPAddr(unixIp uint32, port uint16) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   net.IPv4(byte(unixIp), byte(unixIp>>8), byte(unixIp>>16), byte(unixIp>>24)),
		Port: int(port),
		Zone: "",
	}
}

// ConnectRequestFromBytes A RequestHandler method that parses a bytearray and returns ConnectRequest instance,
// or an error if the bytearray is not a bson representation of ConnectRequest.
func ConnectRequestFromBytes(requestBytes []byte) (*ConnectRequest, error) {
	var connectRequest ConnectRequest
	err := bson.Unmarshal(requestBytes, &connectRequest)
	if err != nil {
		return nil, err
	}
	return &connectRequest, nil
}

// ConnectResponseToBytes A RequestHandler method that serializes ConnectResponse instance to a bytearray.
func ConnectResponseToBytes(response ConnectResponse) ([]byte, error) {
	return bson.Marshal(response)
}
