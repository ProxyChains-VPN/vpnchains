package ipc

// ConnectRequest A wrapper struct for
// int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen)
// syscall arguments.
type ConnectRequest struct {
	SockFd int32  `bson:"sock_fd"`
	Port   uint16 `bson:"port"`
	Ip     int32  `bson:"ip"`
}

// ConnectResponse A wrapper struct for
// int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen)
// syscall return value and errno (TODO).
type ConnectResponse struct {
	ResultCode int32 `bson:"result_code"`
}

// ReadRequest A wrapper struct for
// ssize_t read(int fd, void *buf, size_t count)
// syscall arguments.
type ReadRequest struct {
	Fd          int32  `bson:"fd"`
	BytesToRead uint64 `bson:"bytes_to_read"`
}

// ReadResponse A wrapper struct for
// ssize_t read(int fd, void *buf, size_t count)
// syscall return value and errno (TODO).
// (TODO - int64 -> ssize_t????)
type ReadResponse struct {
	Buffer    []byte `bson:"buffer"`
	BytesRead int64  `bson:"bytes_read"`
}

// WriteRequest A wrapper struct for
// ssize_t write(int fd, const void buf[.count], size_t count)
// syscall arguments.
type WriteRequest struct {
	Fd           int32  `bson:"fd"`
	Buffer       []byte `bson:"buffer"`
	BytesToWrite uint64 `bson:"bytes_to_write"`
}

// WriteResponse A wrapper struct for
// ssize_t write(int fd, const void buf[.count], size_t count)
// syscall return value and errno (TODO).
// (TODO - int64 -> ssize_t????)
type WriteResponse struct {
	BytesWritten int64 `bson:"bytes_written"`
}

//
//// CloseRequest A wrapper struct for
//// int close(int fd)
//// syscall arguments.
//type CloseRequest struct {
//	Fd int32 `bson:"fd"`
//}
//
//// CloseResponse A wrapper struct for
//// int close(int fd)
//// syscall return value and errno (TODO).
//type CloseResponse struct {
//	CloseResult int32 `bson:"close_res"`
//}
