package ipc

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// GetRequestType A function that parses a bson bytearray and returns a string that represents
// the type of the request, or an error if the bytearray is not a bson representation of a request with
// a "call" field.
func GetRequestType(requestBytes []byte) (string, error) {
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
