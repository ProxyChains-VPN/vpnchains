package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"vpnchains/gopkg/ipc"
)

var errorReadResponse = ipc.ReadResponse{[]byte(""), -1}
var errorReadResponseBytes, _ = bson.Marshal(errorReadResponse)

func (handler *RequestHandler) processRead(request *ipc.ReadRequest) (*ipc.ReadResponse, error) {
	log.Println(request)
	buf := make([]byte, request.BytesToRead)
	bytesRead, err := handler.tunnel.Read(request.Fd, buf)
	if err != nil {
		return &errorReadResponse, err
	}

	return &ipc.ReadResponse{buf, bytesRead}, nil
}
