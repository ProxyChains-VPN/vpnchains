package ipc

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func processWrite(request *WriteRequest) (*WriteResponse, error) {
	log.Println(string(request.Buffer))
	return &WriteResponse{100500}, nil
}

func processRead(request *ReadRequest) (*ReadResponse, error) {
	//log.Println(request)
	return &ReadResponse{Buffer: []byte("PRIVET"), BytesRead: 6}, nil
}

func processConnect(request *ConnectRequest) (*ConnectResponse, error) {
	//log.Println(request)
	return &ConnectResponse{0}, nil
}

var errorWriteResponse = WriteResponse{BytesWritten: -1}
var errorWriteResponseBytes, _ = bson.Marshal(errorWriteResponse)

var errorReadResponse = ReadResponse{[]byte(""), 0}
var errorReadResponseBytes, _ = bson.Marshal(errorReadResponse)

var errorConnectResponse = ConnectResponse{ResultCode: -1}
var errorConnectResponseBytes, _ = bson.Marshal(errorConnectResponse)

func HandleRequest(requestBytes []byte) ([]byte, error) {
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
		var writeRequest WriteRequest
		err = bson.Unmarshal(requestBytes, &writeRequest)
		if err != nil {
			return errorWriteResponseBytes, err
		}

		response, err = processWrite(&writeRequest)
	case "read":
		var readRequest ReadRequest
		err = bson.Unmarshal(requestBytes, &readRequest)
		if err != nil {
			return errorReadResponseBytes, err
		}

		response, err = processRead(&readRequest)
	case "connect":
		var connectRequest ConnectRequest
		err = bson.Unmarshal(requestBytes, &connectRequest)
		if err != nil {
			return errorConnectResponseBytes, err
		}

		response, err = processConnect(&connectRequest)
	default:
		return nil, errors.New("wrong format or unknown call")
	}

	if err != nil {
		log.Println(err)
	}

	responseBytes, _ := bson.Marshal(response)
	return responseBytes, err
}
