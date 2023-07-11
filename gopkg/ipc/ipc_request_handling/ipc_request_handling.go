package ipc_request_handling

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"vpnchains/gopkg/ipc"
	"vpnchains/gopkg/vpn"
)

type RequestHandler struct {
	tunnel vpn.Tunnel
}

func NewRequestHandler(tunnel vpn.Tunnel) *RequestHandler {
	return &RequestHandler{tunnel: tunnel}
}

func (handler *RequestHandler) HandleRequest(requestBytes []byte) ([]byte, error) {
	err := bson.Raw(requestBytes).Validate()
	if err != nil {
		return nil, err
	}

	callValue := bson.Raw(requestBytes).Lookup("call")
	call, ok := callValue.StringValueOK()
	if !ok {
		log.Println(callValue.Value)
		return nil, errors.New("call is not string")
	}

	var response interface{}

	switch call {
	case "write":
		var writeRequest ipc.WriteRequest
		err = bson.Unmarshal(requestBytes, &writeRequest)
		if err != nil {
			return errorWriteResponseBytes, err
		}

		response, err = handler.processWrite(&writeRequest)
	case "read":
		var readRequest ipc.ReadRequest
		err = bson.Unmarshal(requestBytes, &readRequest)
		if err != nil {
			return errorReadResponseBytes, err
		}

		response, err = handler.processRead(&readRequest)
	case "connect":
		var connectRequest ipc.ConnectRequest
		err = bson.Unmarshal(requestBytes, &connectRequest)
		if err != nil {
			return errorConnectResponseBytes, err
		}

		response, err = handler.processConnect(&connectRequest)
	default:
		return nil, errors.New("wrong format or unknown call")
	}

	if err != nil {
		log.Println(err)
	}

	responseBytes, _ := bson.Marshal(response) // todo err handling
	return responseBytes, err
}
