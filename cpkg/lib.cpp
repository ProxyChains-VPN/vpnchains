#include "lib.h"
#include <fcntl.h>
#include <dlfcn.h>
#include <netinet/in.h>
#include <sys/stat.h>
#include <sys/un.h>
#include <assert.h>
#include <string.h>
#include <gnu/lib-names.h>
#include <stdbool.h>
#include <libbson-1.0/bson/bson.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <map>

std::map<int, struct sockaddr_in> udp_connections;

/*
 * Gets the IPC port from the environment variable VPNCHAINS_IPC_SERVER_PORT.
 */
int ipc_port = -1;
int get_ipc_port(){
    if (ipc_port == -1) {
        errno = 0;
	    ipc_port = strtoul(getenv(IPC_PORT_ENV_VAR), NULL, 10);
	    if (errno != 0) { // CRINGE CRINGE CRINGE CRINGE
            perror("IPC_PORT_ENV_VAR is not set or set incorrectly");
            exit(1);
        }
	    fprintf(stderr, "\n\n\nipc_port = %d\n\n\n", ipc_port);
    }
    return ipc_port;
}


/*
 * Callbacks.
 */
typedef int (*Connect_callback)(int, const struct sockaddr*, socklen_t);
typedef ssize_t (*Sendto_callback)(int, const void*, size_t, int, const struct sockaddr*, socklen_t);
typedef ssize_t (*Recvfrom_callback)(int, void*, size_t, int, struct sockaddr*, socklen_t*);
typedef ssize_t (*Recvmsg_callback)(int, struct msghdr*, int);
typedef ssize_t (*Sendmsg_callback)(int, const struct msghdr*, int);
typedef ssize_t (*Write_callback)(int, const void*, size_t);
typedef ssize_t (*Read_callback)(int, void*, size_t);

Connect_callback __real_connect = NULL;
Sendto_callback __real_sendto = NULL;
Recvfrom_callback __real_recvfrom = NULL;
Recvmsg_callback __real_recvmsg = NULL;
Sendmsg_callback __real_sendmsg = NULL;
Write_callback __real_write = NULL;
Read_callback __real_read = NULL;

ssize_t real_sendmsg(int sockfd, const struct msghdr *msg, int flags){
    if (__real_sendmsg == NULL) {
        __real_sendmsg = (Sendmsg_callback)dlsym(RTLD_NEXT, "sendmsg");
    }
    return __real_sendmsg(sockfd, msg, flags);
}

int real_connect(int fd, const struct sockaddr* sa, socklen_t len) {
    if (__real_connect == NULL) {
        __real_connect = (Connect_callback)dlsym(RTLD_NEXT, "connect");
    }
    return __real_connect(fd, sa, len);
}

ssize_t real_sendto(int s, const void *msg, size_t len, int flags, const struct sockaddr *to, socklen_t tolen){
    if (__real_sendto == NULL) {
        __real_sendto = (Sendto_callback)dlsym(RTLD_NEXT, "sendto");
    }
    return __real_sendto(s, msg, len, flags, to, tolen);
}

ssize_t real_recvfrom(int s, void *buf, size_t len, int flags, struct sockaddr *from, socklen_t *fromlen){
    if (__real_recvfrom == NULL) {
        __real_recvfrom = (Recvfrom_callback)dlsym(RTLD_NEXT, "recvfrom");
    }
    return __real_recvfrom(s, buf, len, flags, from, fromlen);
}

ssize_t real_recvmsg(int sockfd, struct msghdr *msg, int flags){
    if (__real_recvmsg == NULL) {
        __real_recvmsg = (Recvmsg_callback)dlsym(RTLD_NEXT, "recvmsg");
    }
    return __real_recvmsg(sockfd, msg, flags);
}

ssize_t real_write(int fd, const void *buf, size_t count){
    if (__real_write == NULL) {
        __real_write = (Write_callback)dlsym(RTLD_NEXT, "write");
    }
    return __real_write(fd, buf, count);
}

ssize_t real_read(int fd, void *buf, size_t count){
    if (__real_read == NULL) {
        __real_read = (Read_callback)dlsym(RTLD_NEXT, "read");
    }
    return __real_read(fd, buf, count);
}


/*
 * Socket utils.
 */
bool is_sock(int fd) {
    struct stat statbuf;
    fstat(fd, &statbuf);
    return S_ISSOCK(statbuf.st_mode);
}

int socket_sa_family(int fd) {
    if (!is_sock(fd)) {
        return false;
    }

    struct sockaddr addr;
    socklen_t len = sizeof(addr);
    int returnvalue = getsockname(fd, (struct sockaddr *) &addr, &len);
    if (returnvalue == -1) {
        perror("getsockname() failed");
        return false;
    }
    return addr.sa_family;
}

bool is_internet_socket(int fd) {
    int sa = socket_sa_family(fd);
    return sa == AF_INET || sa == AF_INET6;
}

int socket_type(int fd){
    int socktype = 0;
    socklen_t optlen = sizeof(socktype);

    if(-1 == getsockopt(fd, SOL_SOCKET, SO_TYPE, &socktype, &optlen)){
        perror("getsockopt() failed");
        return -1;
    }

    return socktype;
}

// TODO починить
bool is_localhost(const struct sockaddr *addr){
    assert(addr != NULL);
    assert(addr->sa_family != AF_UNIX);

    if (addr->sa_family == AF_INET) {
        struct sockaddr_in* sin = (struct sockaddr_in*)addr;
        unsigned int ip = sin->sin_addr.s_addr;
        return ip == 0 || ip == 0x0100007f;
    }

    if (addr->sa_family == AF_INET6) {
        struct sockaddr_in6* sin6 = (struct sockaddr_in6*)addr;
        unsigned char* ip = sin6->sin6_addr.s6_addr;

        const char *ip6str = "::1";
        struct in6_addr result;
        assert(inet_pton(AF_INET6, ip6str, &result) == 1);

        if (memcmp(ip, &result, sizeof(struct in6_addr)) == 0) {
            return true;
        }

        const char *ip6str2 = "::";
        struct in6_addr result2;
        assert(inet_pton(AF_INET6, ip6str2, &result2) == 1);

        return memcmp(ip, &result2, sizeof(struct in6_addr)) == 0;
    }

    return false;
}

/*
 * Bson utils.
 */
bool is_valid(const bson_t* bson){
    if (!bson_validate(
            bson,
            bson_validate_flags_t(BSON_VALIDATE_UTF8
            | BSON_VALIDATE_DOLLAR_KEYS
            | BSON_VALIDATE_DOT_KEYS
            | BSON_VALIDATE_UTF8_ALLOW_NULL
            | BSON_VALIDATE_EMPTY_KEYS),
            NULL)) {
        write(2, "Response bson is not valid\n", 27);
        return false;
    }
    return true;
}

/*
 * IPC utils.
 */
int connect_local_socket(int fd) {
    static bool called = false;
    static struct sockaddr_in name;
    if (!called) {
        memset(&name, 0, sizeof(struct sockaddr_in));
        name.sin_family = AF_INET;
        name.sin_port = htons(get_ipc_port());
        if(inet_pton(AF_INET, "127.0.0.1", &name.sin_addr) <= 0){
            perror("inet_pton failed");
            return -1;
	}
        called = true;
    }

    int tmp_sock_connect_res = real_connect(fd, (const struct sockaddr*)&name, sizeof(name));
    if (tmp_sock_connect_res == -1) {
        perror("Connect() tmp socket failed");
        //fprintf(stderr, "%s\n", name.sun_path);
        close(fd);
        return -1;
    }

    return tmp_sock_connect_res;
}

/*
 * Function overrides.
 */
SO_EXPORT int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen) {
    if (addr == NULL) {
        errno = EFAULT;
        return -1;
    }

    if (!is_sock(sock_fd)) {
        errno = ENOTSOCK;
        return -1;
    }

    if (socket_sa_family(sock_fd) == AF_UNIX || is_localhost(addr)) {
        return real_connect(sock_fd, addr, addrlen);
    }

    if (socket_sa_family(sock_fd) == AF_INET6) {
        errno = EAFNOSUPPORT; // todo
        return -1;
    }

    struct sockaddr_in *sin = (struct sockaddr_in*)addr;

    if (socket_type(sock_fd) & SOCK_DGRAM) {
        udp_connections[sock_fd] = *sin;
        return real_connect(sock_fd, addr, addrlen);
    }

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "connect");
    BSON_APPEND_INT32(&bson_request, "sock_fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "port", ntohs(sin->sin_port));
    BSON_APPEND_INT32(&bson_request, "ip", sin->sin_addr.s_addr);

    unsigned int unixIp = sin->sin_addr.s_addr;
    fprintf(stderr, "[line124]connecting to %u.%u.%u.%u:%u\n\n", (unsigned char) unixIp, (unsigned char)(unixIp>>8), (unsigned char)(unixIp>>16), (unsigned char)(unixIp>>24), ntohs(sin->sin_port));

    int flags = fcntl(sock_fd, F_GETFL, 0);
    if (flags & O_NONBLOCK) {
        fcntl(sock_fd, F_SETFL, !O_NONBLOCK); // todo эээээээээ мб не так??????
        if (-1 == connect_local_socket(sock_fd)){
            write(2, "Failed to connect UNIX socket\n", 30);
            return -1;
        }
        fcntl(sock_fd, F_SETFL, O_NONBLOCK);
    } else {
        if (-1 == connect_local_socket(sock_fd)){
            write(2, "Failed to connect UNIX socket\n", 30);
            return -1;
        }
    }

    int bytes_written = write(sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (bytes_written == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    if (!is_valid(bson_response)){
        return -1;
    }

    int res = -1;

    bson_iter_t iter;
    bson_iter_t result_code;
    if (!bson_iter_init(&iter, bson_response)){
        perror("Failed to parse bson: bson_iter_init");
        return -1;
    }

    if (!bson_iter_find_descendant(&iter, "result_code", &result_code)){
        perror("Failed to parse bson: can't find 'result_code'");
        return -1;
    }

    if (!BSON_ITER_HOLDS_INT32(&result_code)){
        perror("Failed to parse bson: 'result_code' is not int32");
        return -1;
    }

    res = bson_iter_int32(&result_code);

    bson_reader_destroy(reader);

    fprintf(stderr, "\n[line 172] connect result %d\n\n", res);
    return res;
}

SO_EXPORT ssize_t sendto(int s, const void *msg, size_t len, int flags, const struct sockaddr *to, socklen_t tolen) {
    if (!is_sock(s)) {
        errno = ENOTSOCK;
        return -1;
    }

    if (socket_sa_family(s) == AF_INET6) {
        fprintf(stderr, "sendto: AF_INET6 not supported\n");
        return real_sendto(s, msg, len, flags, to, tolen);
//        errno = EAFNOSUPPORT;
//        return -1;
    }

    else if (socket_sa_family(s) == AF_UNIX) {
        return real_sendto(s, msg, len, flags, to, tolen);
    }

    else if (socket_sa_family(s) == AF_INET && (socket_type(s) & SOCK_DGRAM) && (to == NULL || !is_localhost(to))) {
        if (to == NULL && udp_connections.find(s) == udp_connections.end()) {
            fprintf(stderr, "sendto: to is NULL\n");
            fprintf(stderr, "ыутвещ not local");
            errno = ECONNREFUSED;
            return -1;
        }

        struct sockaddr_in* sin;
        if (to == NULL) {
            fprintf(stderr, "sendto: got TO from map\n");
            sin = &udp_connections[s];
        }
        else {
            sin = (struct sockaddr_in*)to;
        }
//
//
//
        int ipc_sock_fd = socket(AF_INET, SOCK_DGRAM, 0);
        if (ipc_sock_fd == -1) {
            perror("Failed to open udp socket");
            return -1;
        }
        int sock_flags = fcntl(s, F_GETFL, 0);
        fcntl(ipc_sock_fd, F_SETFL, sock_flags);

        struct sockaddr_in name;
        memset(&name, 0, sizeof(struct sockaddr_in));
        name.sin_family = AF_INET;
        name.sin_port = htons(get_ipc_port());
        if(inet_pton(AF_INET, "127.0.0.1", &name.sin_addr) <= 0){
            perror("inet_pton failed");
            return -1;
        }


        bson_t bson_request = BSON_INITIALIZER;
        BSON_APPEND_UTF8(&bson_request, "call", "sendto");
        BSON_APPEND_BINARY(&bson_request, "msg", BSON_SUBTYPE_BINARY, (const unsigned char*) msg, len);
        BSON_APPEND_INT64(&bson_request, "msg_len", len);


        BSON_APPEND_INT32(&bson_request, "dest_ip", sin->sin_addr.s_addr);
        BSON_APPEND_INT32(&bson_request, "dest_port", ntohs(sin->sin_port));
        BSON_APPEND_INT64(&bson_request, "pid", getpid());
        BSON_APPEND_INT32(&bson_request, "fd", s);

        real_sendto(ipc_sock_fd, bson_get_data(&bson_request), bson_request.len, 0, (const struct sockaddr*)&name, sizeof(name));
//
        fprintf(stderr, "sendto: sent request\n");
//
//
        bson_destroy(&bson_request);
//        real_sendto(s, msg, len, flags, to, tolen);
        return len;
    }

    else {
        return real_sendto(s, msg, len, flags, to, tolen);
    }
}

SO_EXPORT ssize_t send(int s, const void *msg, size_t len, int flags) {
    if (!is_sock(s)) {
        errno = ENOTSOCK;
        return -1;
    }
    return sendto(s, msg, len, flags, NULL, 0);
}

SO_EXPORT ssize_t recv(int s, void *buf, size_t len, int flags) {
    if (!is_sock(s)) {
        errno = ENOTSOCK;
        return -1;
    }
    return recvfrom(s, buf, len, flags, NULL, NULL);
}

SO_EXPORT ssize_t read(int fd, void *buf, size_t count) {
    if (is_sock(fd)) {
        return recv(fd, buf, count, 0);
    } else {
        return real_read(fd, buf, count);
    }
}

SO_EXPORT ssize_t write(int fd, const void *buf, size_t count) {
    if (is_sock(fd)) {
        return send(fd, buf, count, 0);
    } else {
        return real_write(fd, buf, count);
    }
}

SO_EXPORT ssize_t sendmsg(int s, const struct msghdr *msg, int flags) {
//    fprintf(stderr, "sendmsg\n");
    if (!is_sock(s)) {
        errno = ENOTSOCK;
        return -1;
    }

    if (socket_sa_family(s) == AF_INET && (socket_type(s) & SOCK_DGRAM)) {
        return sendto(s, msg->msg_iov[0].iov_base, msg->msg_iov[0].iov_len, flags, (const sockaddr*) msg->msg_name, msg->msg_namelen);
    }

    return real_sendmsg(s, msg, flags);
}

SO_EXPORT ssize_t recvfrom(int s, void *buf, size_t len, int flags, struct sockaddr *from, socklen_t *fromlen){
    if (!is_sock(s)) {
        errno = ENOTSOCK;
        return -1;
    }

    if (socket_sa_family(s) == AF_INET6) {
        fprintf(stderr, "recvfrom: AF_INET6 not supported\n");
//        errno = EAFNOSUPPORT;
//        return -1;
        return real_recvfrom(s, buf, len, flags, from, fromlen);
    }

    if (socket_sa_family(s) == AF_UNIX) {
        return real_recvfrom(s, buf, len, flags, from, fromlen);
    }

    if (socket_sa_family(s) == AF_INET && (socket_type(s) & SOCK_DGRAM)) {
        fprintf(stderr, "IN SO recvfrom fd %d pid %d\n", s, getpid());

//        return real_recvfrom(s, buf, len, flags, from, fromlen);

        int ipc_sock_fd = socket(AF_INET, SOCK_DGRAM, 0);
        if (ipc_sock_fd == -1) {
            perror("Failed to open udp socket");
            return -1;
        }

        struct sockaddr_in name;
        memset(&name, 0, sizeof(struct sockaddr_in));
        name.sin_family = AF_INET;
        name.sin_port = htons(get_ipc_port());
        if(inet_pton(AF_INET, "127.0.0.1", &name.sin_addr) <= 0){
            perror("inet_pton failed");
            return -1;
        }

        bson_t bson_request = BSON_INITIALIZER;
        BSON_APPEND_UTF8(&bson_request, "call", "recvfrom");
        BSON_APPEND_INT64(&bson_request, "pid", getpid());
        BSON_APPEND_INT32(&bson_request, "fd", s);
        BSON_APPEND_INT64(&bson_request, "msg_len", len);


        real_sendto(ipc_sock_fd, bson_get_data(&bson_request), bson_request.len, 0, (const struct sockaddr*)&name, sizeof(name));

        fprintf(stderr, "recvfrom: sent request\n");

        bson_destroy(&bson_request);

	    uint8_t *buf = (uint8_t*)malloc(len+1024);
        socklen_t name_len = sizeof(name);


        int bb_read;

        int flags = fcntl(ipc_sock_fd, F_GETFL, 0);
        if (flags & O_NONBLOCK) {
            fcntl(ipc_sock_fd, F_SETFL, flags & ~O_NONBLOCK);
            bb_read = real_recvfrom(ipc_sock_fd, (void*)buf, len+1024, 0, (struct sockaddr*)&name, &name_len);
            if (-1 == bb_read){
                perror("recvfrom() ipc socket failed:\n");
                return -1;
            }
            fcntl(ipc_sock_fd, F_SETFL, flags | O_NONBLOCK);
        } else {
             bb_read = real_recvfrom(ipc_sock_fd, (void*)buf, len+1024, 0, (struct sockaddr*)&name, &name_len);
             if (-1 == bb_read){
                 perror("recvfrom() ipc socket failed:\n");
                 return -1;
             }
        }

        fprintf(stderr, "RECVFROM: got response, bytes read %d\n", bb_read);

        bson_reader_t* reader = bson_reader_new_from_data(buf, bb_read);

        const bson_t* bson_response = bson_reader_read(reader, NULL);
        if (!is_valid(bson_response)){
            return -1;
        }

        int res = -1;

        bson_iter_t iter;
        bson_iter_t bson_bytes_read;
        bson_iter_t bson_msg;
        bson_iter_t bson_src_ip;
        bson_iter_t bson_src_port;
        ssize_t bytes_read;
        int src_ip;
        int src_port;
        void *binary_data;
        if (!bson_iter_init(&iter, bson_response)){
            perror("Failed to parse bson: bson_iter_init");
            return -1;
        }

        if (!bson_iter_find_descendant(&iter, "bytes_read", &bson_bytes_read)){
            perror("Failed to parse bson: can't find 'bytes_read'");
            return -1;
        }

        if (!bson_iter_find_descendant(&iter, "msg", &bson_msg)){
            perror("Failed to parse bson: can't find 'msg'");
            return -1;
        }

        if (!bson_iter_find_descendant(&iter, "src_ip", &bson_src_ip)){
            perror("Failed to parse bson: can't find 'src_ip'");
            return -1;
        }

        if (!bson_iter_find_descendant(&iter, "src_port", &bson_src_port)){
            perror("Failed to parse bson: can't find 'src_port'");
            return -1;
        }

        if (!BSON_ITER_HOLDS_INT64(&bson_bytes_read)){
            perror("Failed to parse bson: 'bytes_read' is not int64");
            return -1;
        }

        if (!BSON_ITER_HOLDS_BINARY(&bson_msg)){
             perror("Failed to parse bson: 'msg' is not binary");
            return -1;
        }

        if (!BSON_ITER_HOLDS_INT32(&bson_src_ip)){
            perror("Failed to parse bson: 'src_ip' is not int32");
            return -1;
        }

        if (!BSON_ITER_HOLDS_INT32(&bson_src_port)){
            perror("Failed to parse bson: 'src_port' is not int32");
            return -1;
        }

        bytes_read = bson_iter_int64(&bson_bytes_read);

        if (bytes_read == -1) {
            return -1;
        }
        bson_iter_binary(&bson_msg, NULL, (unsigned int *) &bytes_read, (const uint8_t**)&binary_data);
        memcpy(buf, binary_data, std::min(size_t(bytes_read), len));
        src_ip = bson_iter_int32(&bson_src_ip);
        src_port = bson_iter_int32(&bson_src_port);

        if (from != NULL){
            struct sockaddr_in* from_in = (struct sockaddr_in*)from;
            from_in->sin_family = AF_INET;
            from_in->sin_addr.s_addr = src_ip;
            from_in->sin_port = src_port;
            *fromlen = sizeof(from_in);
        }

        bson_reader_destroy(reader);

        if (-1 == close(ipc_sock_fd)){
            perror("Close() ipc socket failed");
            return -1;
        }

        return bytes_read;
    }

    return real_recvfrom(s, buf, len, flags, from, fromlen);
}

SO_EXPORT ssize_t recvmsg(int s, struct msghdr *msg, int flags) { // todo
    if (!is_sock(s)) {
        fprintf(stderr, "not a socket\n");
        errno = ENOTSOCK;
        return -1;
    }

//    if (socket_sa_family(s) != AF_INET6) {
//        fprintf(stderr, "AF INET 6 UNSUPPORTED\n");
//        errno = ECONNREFUSED;
//        return -1;
//    }

//    if (socket_sa_family(s) == AF_UNIX || socket_sa_family(s) == AF_INET6) { // todo ??? ok???
//        return real_recvmsg(s, msg, flags);
//    }

    if (socket_sa_family(s) == AF_INET && (socket_type(s) & SOCK_DGRAM)) {
        fprintf(stderr, "recvmsg\n fd %d pid %d\n", s, getpid());

        int retval = recvfrom(s, msg->msg_iov[0].iov_base, msg->msg_iov[0].iov_len, flags, NULL, NULL);

//        struct sockaddr_in name;
//        memset(&name, 0, sizeof(struct sockaddr_in));
//        name.sin_family = AF_INET;
//        name.sin_port = htons(get_ipc_port());
//        if(inet_pton(AF_INET, "127.0.0.1", &name.sin_addr) <= 0){
//            perror("inet_pton failed");
//            return -1;
//        }
        msg->msg_name = NULL;
    }

    return real_recvmsg(s, msg, flags);
//    fprintf(stderr, "recvmsg\n fd %d pid %d\n", sockfd, getpid());
//    return recvfrom(sockfd, msg->msg_iov[0].iov_base, msg->msg_iov[0].iov_len, flags, NULL, NULL);
}