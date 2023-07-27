package udp_ipc_request

import (
	"go.mongodb.org/mongo-driver/bson"
	"net"
)

// ErrorRecvfromResponse A response that is sent to caller process if the sendto syscall was unsuccessful
// (if there were no packet recieved from tunnel side).
var ErrorRecvfromResponse = RecvfromResponse{
	BytesRead: -1,
	Msg:       []byte{},
}

// UnixIpPortToUDPAddr A function that converts POSIX-style IP address (that is an unsigned 64-bit value)
// and port (that is an unsigned 16-bit value) to a net.UDPAddr instance.
// unixIp - POSIX-style IP address.
// port - port.
func UnixIpPortToUDPAddr(unixIp uint32, port uint16) *net.UDPAddr {
	return &net.UDPAddr{
		IP:   net.IPv4(byte(unixIp), byte(unixIp>>8), byte(unixIp>>16), byte(unixIp>>24)),
		Port: int(port),
		Zone: "",
	}
}

// RecvfromRequestFromBytes A function that parses a bytearray and returns RecvfromRequest instance,
// or an error if the bytearray is not a bson representation of RecvfromRequest.
// requestBytes - bytearray to parse.
func RecvfromRequestFromBytes(requestBytes []byte) (*RecvfromRequest, error) {
	var recvfromRequest RecvfromRequest
	err := bson.Unmarshal(requestBytes, &recvfromRequest)
	if err != nil {
		return nil, err
	}
	return &recvfromRequest, nil
}

// RecvfromResponseToBytes A function that serializes RecvfromResponse instance to a bytearray.
// response - RecvfromResponse instance to serialize.
func RecvfromResponseToBytes(response RecvfromResponse) ([]byte, error) {
	return bson.Marshal(response)
}

// SendtoRequestFromBytes A function that parses a bytearray and returns SendtoRequest instance,
// or an error if the bytearray is not a bson representation of SendtoRequest.
// requestBytes - bytearray to parse.
func SendtoRequestFromBytes(requestBytes []byte) (*SendtoRequest, error) {
	var sendtoRequest SendtoRequest
	err := bson.Unmarshal(requestBytes, &sendtoRequest)
	if err != nil {
		return nil, err
	}
	return &sendtoRequest, nil
}
