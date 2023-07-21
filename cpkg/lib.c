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

int ipc_port = -1;
int get_ipc_port(){
    if (ipc_port == -1) {
        errno = 0;
	    ipc_port = strtoul(getenv(IPC_PORT), NULL, 10);
	    if (errno != 0) { // CRINGE CRINGE CRINGE CRINGE
            perror("IPC_PORT is not set or set incorrectly");
            exit(1);
        }
	    fprintf(stderr, "\n\n\nipc_port = %d\n\n\n", ipc_port);
    }
    return ipc_port;
}

//unsigned int local_network_mask[4] = { 10, 127, 4268, 43200 };
//10.0.0.0/8, 127.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16

typedef int (*Connect_callback)(int, const struct sockaddr*, socklen_t);
typedef ssize_t (*Sendto_callback)(int, const void*, size_t, int, const struct sockaddr*, socklen_t);
typedef ssize_t (*Recvfrom_callback)(int, void*, size_t, int, struct sockaddr*, socklen_t*);

Connect_callback __real_connect = NULL;
Sendto_callback __real_sendto = NULL;
Recvfrom_callback __real_recvfrom = NULL;

int real_connect(int fd, const struct sockaddr* sa, socklen_t len) {
    if (__real_connect == NULL) {
        void *h_dl = RTLD_NEXT;
        if (h_dl == NULL) {
            exit(66);
        }

        __real_connect = (Connect_callback)dlsym(h_dl, "connect");
    }
    return __real_connect(fd, sa, len);
}

ssize_t real_sendto(int s, const void *msg, size_t len, int flags, const struct sockaddr *to, socklen_t tolen){
    if (__real_sendto == NULL) {
        void *h_dl = RTLD_NEXT;
        if (h_dl == NULL) {
            exit(66);
        }

        __real_sendto = (Sendto_callback)dlsym(h_dl, "sendto");
    }
    return __real_sendto(s, msg, len, flags, to, tolen);
}

ssize_t real_recvfrom(int s, void *buf, size_t len, int flags, struct sockaddr *from, socklen_t *fromlen){
    if (__real_recvfrom == NULL) {
        void *h_dl = RTLD_NEXT;
        if (h_dl == NULL) {
            exit(66);
        }

        __real_recvfrom = (Recvfrom_callback)dlsym(h_dl, "recvfrom");
    }
    return __real_recvfrom(s, buf, len, flags, from, fromlen);
}

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

bool is_unix_socket(int fd) {
    return socket_sa_family(fd) == AF_UNIX;
}

bool is_ipv4_socket(int fd) {
    return socket_sa_family(fd) == AF_INET;
}

bool is_ipv6_socket(int fd) {
    return socket_sa_family(fd) == AF_INET6;
}

bool is_internet_socket(int fd) {
    return is_ipv4_socket(fd) || is_ipv6_socket(fd);
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

bool is_stream_socket(int fd){
    int socktype = 0;
    socklen_t optlen = sizeof(socktype);

    if(-1 == getsockopt(fd, SOL_SOCKET, SO_TYPE, &socktype, &optlen)){
        perror("getsockopt() failed");
        return false;
    }

    return socktype & SOCK_STREAM ? true : false;
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

//    for(int i; i < 4; i++){
//        if(((ip & local_network_mask[i]) ^ local_network_mask[i]) == 0){
//            return true;
//        }
//    }


}

bool is_valid(const bson_t* bson){
    if (!bson_validate(
            bson,
            BSON_VALIDATE_UTF8
            | BSON_VALIDATE_DOLLAR_KEYS
            | BSON_VALIDATE_DOT_KEYS
            | BSON_VALIDATE_UTF8_ALLOW_NULL
            | BSON_VALIDATE_EMPTY_KEYS,
            NULL)) {
        write(2, "Response bson is not valid\n", 27);
        return false;
    }
    return true;
}


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

SO_EXPORT int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen) {
    if (addr == NULL) {
        errno = EFAULT;
        return -1;
    }

    if ((socket_type(sock_fd) & SOCK_DGRAM) && !is_localhost(addr)) {
        errno = ECONNREFUSED;
        return -1;
    }

    if (is_ipv6_socket(sock_fd) && !is_localhost(addr)) {
        errno = EAFNOSUPPORT; // todo
        return -1;
    }

    if (is_unix_socket(sock_fd) || is_localhost(addr)) {
        return real_connect(sock_fd, addr, addrlen);
    }

    struct sockaddr_in* sin = (struct sockaddr_in*)addr;

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "connect");
    BSON_APPEND_INT32(&bson_request, "sock_fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "port", ntohs(sin->sin_port));
    BSON_APPEND_INT32(&bson_request, "ip", sin->sin_addr.s_addr);

    unsigned int unixIp = sin->sin_addr.s_addr;
    fprintf(stderr, "[line124]connecting to %u.%u.%u.%u:%u\n\n", (unsigned char) unixIp, (unsigned char)(unixIp>>8), (unsigned char)(unixIp>>16), (unsigned char)(unixIp>>24), ntohs(sin->sin_port));

    int flags = fcntl(sock_fd, F_GETFL, 0);
    if (flags & O_NONBLOCK) {
        fcntl(sock_fd, F_SETFL, !O_NONBLOCK);
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

SO_EXPORT ssize_t sendto(int s, const void *msg, size_t len, int flags, const struct sockaddr *to, socklen_t tolen){
    if (is_internet_socket(s) && (socket_type(s) & SOCK_DGRAM) && (to == NULL || !is_localhost(to))) {
        fprintf(stderr, "ыутвещ not local");
        errno = ECONNREFUSED;
        return -1;
    }
    return real_sendto(s, msg, len, flags, to, tolen);
}

SO_EXPORT ssize_t recvfrom(int s, void *buf, size_t len, int flags, struct sockaddr *from, socklen_t *fromlen){
    if (is_internet_socket(s) && (socket_type(s) & SOCK_DGRAM) && (from == NULL && !is_localhost(from))) {
        fprintf(stderr, "recv not local");
        errno = ECONNREFUSED;
        return -1;
    }

    return real_recvfrom(s, buf, len, flags, from, fromlen);
}

//SO_EXPORT int close(int fd) {
//    fprintf(stderr, "closing fd %d", fd);
//    return shutdown(fd, SHUT_RDWR);
//}
