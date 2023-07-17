#include "lib.h"
#include <fcntl.h>
#include <dlfcn.h>
#include <netinet/in.h>
#include <sys/stat.h>
#include <sys/un.h>
#include <assert.h>
#include <string.h>
#include <stdio.h>
#include <gnu/lib-names.h>
#include <stdbool.h>
#include <libbson-1.0/bson/bson.h>
#include <pthread.h>

#include <stdlib.h>
#include <stdio.h>

unsigned int local_network_mask[5] = { 10, 127, 4268, 43200, 0 };
//10.0.0.0/8, 127.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, 0.0.0.0

typedef int (*Connect_callback)(int, const struct sockaddr*, socklen_t);

Connect_callback __real_connect = NULL;

int real_connect(int fd, const struct sockaddr* sa, socklen_t len) {
    if (__real_connect == NULL) {
        void *h_dl = RTLD_NEXT;
        if (h_dl == NULL) {
            exit(66);
        }

        __real_connect = (Connect_callback)dlsym(h_dl, "connect");
    }
    __real_connect(fd, sa, len);
}

bool is_internet_socket(int fd) {
    struct stat statbuf;
    fstat(fd, &statbuf);
    if (!S_ISSOCK(statbuf.st_mode)){
        return false;
    }

    struct sockaddr addr;
    socklen_t len = sizeof(addr);
    int returnvalue = getsockname(fd, (struct sockaddr *) &addr, &len);
    if (returnvalue == -1) {
        perror("getsockname() failed");
        return false;
    }
    return addr.sa_family == AF_INET6 || addr.sa_family == AF_INET;
}

bool is_stream_socket(int fd){
    int socktype = 0;
    socklen_t optlen = sizeof(socktype);

    if(-1 == getsockopt(fd, SOL_SOCKET, SO_TYPE, &socktype, &optlen)){
        perror("getsockopt() failed");
        return false;
    }

    return socktype == SOCK_STREAM;
}

// TODO починить
bool is_localhost(const struct sockaddr *addr){
    struct sockaddr_in* sin = (struct sockaddr_in*)addr;
    unsigned int ip = sin->sin_addr.s_addr;

    for(int i; i < 5; i++){
        if(((ip & local_network_mask[i]) ^ local_network_mask[i]) == 0){
            return true;
        }
    }
    return false;
//    return ip == 16777343 || ip == 0;
}

bool is_valid(const bson_t* bson);

int connect_local_socket(int fd) {
    static bool called = false;
    static struct sockaddr_in name;
    if (!called) {
        memset(&name, 0, sizeof(struct sockaddr_in));
        name.sin_family = AF_INET;
        name.sin_port = htons(45454);
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

pthread_mutex_t lock;

SO_EXPORT int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen) {
    if (!is_internet_socket(sock_fd) || !is_stream_socket(sock_fd) || is_localhost(addr)) {
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

//    return real_connect(sock_fd, addr, addrlen);

//    pthread_mutex_lock(&lock);

    if(-1 == connect_local_socket(sock_fd)){
//    pthread_mutex_unlock(&lock); 
        write(2, "Failed to connect UNIX socket\n", 30);
        return -1;
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

//SO_EXPORT int close(int fd) {
//    fprintf(stderr, "closing fd %d", fd);
//    return shutdown(fd, SHUT_RDWR);
//}

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
