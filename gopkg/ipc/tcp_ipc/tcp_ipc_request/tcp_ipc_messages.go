package tcp_ipc_request

// ConnectRequest A wrapper struct for
// int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen)
// syscall arguments.
type ConnectRequest struct {
	SockFd int32  `bson:"sock_fd"`
	Ip     int32  `bson:"ip"`
	Port   uint16 `bson:"port"`
}

// ConnectResponse A wrapper struct for
// int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen)
// syscall return value and errno (TODO).
type ConnectResponse struct {
	ResultCode int32 `bson:"result_code"`
}
