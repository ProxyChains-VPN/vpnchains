package ipc_request

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"vpnchains/gopkg/vpn"
)

type RequestHandler struct {
	tunnel vpn.Tunnel
}

func NewRequestHandler(tunnel vpn.Tunnel) *RequestHandler {
	return &RequestHandler{tunnel: tunnel}
}

func (handler *RequestHandler) GetRequestType(requestBytes []byte) (string, error) {
	err := bson.Raw(requestBytes).Validate()
	if err != nil {
		return "", err
	}

	callValue := bson.Raw(requestBytes).Lookup("call")
	call, ok := callValue.StringValueOK()
	if !ok {
		log.Println(callValue.Value)
		return "", errors.New("call is not string")
	} else {
		return call, nil
	}
}
