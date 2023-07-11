package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"vpnchains/gopkg/ipc"
)

var errorReadResponse = ipc.ReadResponse{[]byte(""), 0}
var errorReadResponseBytes, _ = bson.Marshal(errorReadResponse)

func (handler *RequestHandler) processRead(request *ipc.ReadRequest) (*ipc.ReadResponse, error) {
	//log.Println(request)
	return &ipc.ReadResponse{Buffer: []byte("PRIVET"), BytesRead: 6}, nil
}
