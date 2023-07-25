package udp_ipc_request

// SendtoRequest A wrapper struct for
// ssize_t sendto(int sockfd, const void *buf, size_t len, int flags,
//
//	const struct sockaddr *dest_addr, socklen_t addrlen);
//
// syscall arguments.
type SendtoRequest struct {
	Msg      []byte `bson:"msg"`
	MsgLen   uint64 `bson:"msg_len"`
	DestIp   int32  `bson:"dest_ip"`
	DestPort int16  `bson:"dest_port"`
	Pid      int64  `bson:"pid"`
	Fd       int32  `bson:"fd"`
}

// RecvfromRequest A wrapper struct for
// ssize_t recvfrom(int sockfd, void *buf, size_t len, int flags,
//
//	struct sockaddr *src_addr, socklen_t *addrlen);
//
// syscall arguments.
type RecvfromRequest struct {
	Pid     int64  `bson:"pid"`
	Fd      int32  `bson:"fd"`
	MsgLen  uint64 `bson:"msg_len"`
	SrcIp   int32  `bson:"src_ip"`
	SrcPort int16  `bson:"src_port"`
}

// RecvfromResponse A wrapper struct for
// ssize_t recvfrom(int sockfd, void *buf, size_t len, int flags,
//
//	struct sockaddr *src_addr, socklen_t *addrlen);
//
// syscall return value and errno (TODO).
type RecvfromResponse struct {
	BytesRead int64  `bson:"bytes_read"`
	Msg       []byte `bson:"msg"`
}
