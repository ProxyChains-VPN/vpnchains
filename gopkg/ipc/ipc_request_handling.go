package ipc

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func processWrite(request *WriteRequest) (*WriteResponse, error) {
	log.Println(request)
	return &WriteResponse{100500}, nil
}

func processRead(request *ReadRequest) (*ReadResponse, error) {
	log.Println(request)
	return &ReadResponse{"animeaboba", 666}, nil
}

func processConnect(request *ConnectRequest) (*ConnectResponse, error) {
	log.Println(request)
	return &ConnectResponse{-1}, nil
}

func HandleRequest(request []byte) ([]byte, error) {
	err := bson.Raw(request).Validate()
	if err != nil {
		return nil, err
	}

	callValue := bson.Raw(request).Lookup("Call")
	call, ok := callValue.StringValueOK()
	if !ok {
		log.Println(callValue.Value)
		return nil, errors.New("call is not string")
	}

	var response interface{}

	switch call {
	case "write":
		var writeRequest WriteRequest
		err = bson.Unmarshal(request, &writeRequest)
		if err != nil {
			return nil, err
		}

		response, err = processWrite(&writeRequest)
	case "read":
		var readRequest ReadRequest
		err = bson.Unmarshal(request, &readRequest)
		if err != nil {
			return nil, err
		}

		response, err = processRead(&readRequest)
	case "connect":
		var connectRequest ConnectRequest
		err = bson.Unmarshal(request, &connectRequest)
		if err != nil {
			return nil, err
		}

		response, err = processConnect(&connectRequest)
	default:
		return nil, errors.New("wrong format or unknown call")
	}

	if err != nil {
		return nil, err
	}

	return bson.Marshal(response)
}

func main() {

}
