package ipc_request_handling

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"vpnchains/gopkg/ipc"
)

var errorWriteResponse = ipc.WriteResponse{BytesWritten: -1}
var errorWriteResponseBytes, _ = bson.Marshal(errorWriteResponse)

func (handler *RequestHandler) processWrite(request *ipc.WriteRequest) (*ipc.WriteResponse, error) {
	log.Println("write", request)
	bytesWritten, err := handler.tunnel.Write(request.Fd, request.Buffer[:request.BytesToWrite])
	log.Println("BYTES WRITTEN", bytesWritten)
	if err != nil {
		return &errorWriteResponse, err
	}
	return &ipc.WriteResponse{bytesWritten}, nil
}
