package ipc

type ConnectRequest struct {
	SockFd int32  `bson:"sock_fd"`
	Port   uint16 `bson:"port"`
	Ip     uint32 `bson:"ip"`
}

type ConnectResponse struct {
	ResultCode int32 `bson:"result_code"`
}

type ReadRequest struct {
	Fd          int32 `bson:"fd"`
	BytesToRead int32 `bson:"bytes_to_read"`
}

type ReadResponse struct {
	Buffer    string `bson:"buffer"`
	BytesRead int32  `bson:"bytes_read"`
}

type WriteRequest struct {
	Fd           int32  `bson:"fd"`
	Buffer       []byte `bson:"buffer"`
	BytesToWrite int32  `bson:"bytes_to_write"`
}

type WriteResponse struct {
	BytesWritten int32 `bson:"bytes_written"`
}

type CloseRequest struct {
	Fd int32 `bson:"fd"`
}

type CloseResponse struct {
	CloseResult int32 `bson:"close_res"`
}
