package ipc

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func processWrite(request *WriteRequest) (*WriteResponse, error) {
	log.Fatalln(request)
	return nil, nil
}

func processRead(request *ReadRequest) (*ReadResponse, error) {
	log.Fatalln(request)
	return nil, nil
}

func processConnect(request *ConnectRequest) (*ConnectResponse, error) {
	log.Fatalln(request)
	return nil, nil
}

func HandleRequest(request []byte) ([]byte, error) {
	err := bson.Raw(request).Validate()
	if err != nil {
		return nil, err
	}
	call := bson.Raw(request).Lookup("call").StringValue()

	switch call {
	case "write":
		var writeRequest WriteRequest
		err = bson.Unmarshal(request, &writeRequest)
		if err != nil {
			return nil, err
		}

		writeResponse, err := processWrite(&writeRequest)
		if err != nil {
			return nil, err
		}
		return bson.Marshal(writeResponse)
	case "read":
		var readRequest ReadRequest
		err = bson.Unmarshal(request, &readRequest)
		if err != nil {
			return nil, err
		}

		readResponse, err := processRead(&readRequest)
		if err != nil {
			return nil, err
		}
		return bson.Marshal(readResponse)
	case "connect":
		var connectRequest ConnectRequest
		err = bson.Unmarshal(request, &connectRequest)
		if err != nil {
			return nil, err
		}

		connectResponse, err := processConnect(&connectRequest)
		if err != nil {
			return nil, err
		}
		return bson.Marshal(connectResponse)
	default:
		return nil, errors.New("wrong format or unknown call")
	} // TODO если кастануть к interface{} сломается???? а то выглядит очень очень плохо
}

func main() {

}
