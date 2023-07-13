package ipc_request_handling

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"vpnchains/gopkg/so_ipc"
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

	switch call {
	case "connect":
		var connectRequest so_ipc.ConnectRequest
		err = bson.Unmarshal(requestBytes, &connectRequest)
		if err != nil {
			return errorConnectResponseBytes, err
		}

		response, err := handler.processConnect(&connectRequest)
		if err != nil {
			return errorConnectResponseBytes, err
		}

		responseBytes, err := bson.Marshal(response) // todo err handling
		if err != nil {
			return errorConnectResponseBytes, err
		}
		return responseBytes, nil
	default:
		return nil, errors.New("wrong format or unknown call")
	}
}
