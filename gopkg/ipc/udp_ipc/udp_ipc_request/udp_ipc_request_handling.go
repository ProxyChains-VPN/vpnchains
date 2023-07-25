package udp_ipc_request

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// RequestHandler A struct that contains fields required for handling requests and responses.
// Actually, it is empty now.
type RequestHandler struct {
}

// NewRequestHandler A function that creates a new RequestHandler instance and returns a pointer to it.
// tunnel - a vpn.Tunnel instance.
func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

// GetRequestType A RequestHandler method that parses a bytearray and returns a string that represents
// the type of the request, or an error if the bytearray is not a bson representation of a request with
// a "call" field.
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
