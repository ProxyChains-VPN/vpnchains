package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"vpnchains/gopkg/ipc"
)

var errorWriteResponse = ipc.WriteResponse{BytesWritten: -1}
var errorWriteResponseBytes, _ = bson.Marshal(errorWriteResponse)

func (handler *RequestHandler) processWrite(request *ipc.WriteRequest) (*ipc.WriteResponse, error) {
	log.Println(string(request.Buffer))
	return &ipc.WriteResponse{100500}, nil
}
