package so_ipc

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
