package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"vpnchains/gopkg/ipc"
)

var errorConnectResponse = ipc.ConnectResponse{ResultCode: -1}
var errorConnectResponseBytes, _ = bson.Marshal(errorConnectResponse)

func (handler *RequestHandler) processConnect(request *ipc.ConnectRequest) (*ipc.ConnectResponse, error) {
	//log.Println(request)
	return &ipc.ConnectResponse{0}, nil
}
