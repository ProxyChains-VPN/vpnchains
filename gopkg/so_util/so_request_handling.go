package so_util

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

func processWrite(request *WriteRequest) (*WriteResponse, error) {
	return nil, nil
}

func processRead(request *ReadRequest) (*ReadResponse, error) {
	return nil, nil
}

func processConnect(request *ConnectRequest) (*ConnectResponse, error) {
	return nil, nil
}

func HandleRequest(request []byte) ([]byte, error) {
	err := bson.Raw(request).Validate()
	if err != nil {
		return nil, err
	}
	call := bson.Raw(request).Lookup().String()

	switch call {
	case "write":
		var writeRequest WriteRequest
		err = bson.Unmarshal(request, &writeRequest)
		writeResponse, err := processWrite(&writeRequest)
		if err != nil {
			return nil, err
		}
		return bson.Marshal(writeResponse)
	case "read":
		var readRequest ReadRequest
		err = bson.Unmarshal(request, &readRequest)
		readResponse, err := processRead(&readRequest)
		if err != nil {
			return nil, err
		}
		return bson.Marshal(readResponse)
	case "connect":
		var connectRequest ConnectRequest
		err = bson.Unmarshal(request, &connectRequest)
		connectResponse, err := processConnect(&connectRequest)
		if err != nil {
			return nil, err
		}
		return bson.Marshal(connectResponse)
	default:
		return nil, errors.New("wrong format or unknown call")
	}
}

func main() {

}
