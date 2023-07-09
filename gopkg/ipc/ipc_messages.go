package ipc

type ConnectRequest struct {
	SockFd int32
	Port   uint16
	Ip     uint32
}

// TODO написать чтобы буковки сериализовывались с сохранением регистра
type ConnectResponse struct {
	ResultCode int32 // 0 или -1
}

type ReadRequest struct {
	Fd          int32
	BytesToRead int32
}

type ReadResponse struct {
	Buffer    string
	BytesRead int32
}

type WriteRequest struct {
	Fd           int32
	Buffer       []byte
	BytesToWrite int32
}

type WriteResponse struct {
	BytesWritten int32
}
