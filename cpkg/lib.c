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

#include <stdlib.h>
#include <stdio.h>

typedef ssize_t (*Read_callback)(int, void*, size_t);
typedef int (*Write_callback)(int, const void*, size_t);
typedef int (*Connect_callback)(int, const struct sockaddr*, socklen_t);
typedef int (*Close_callback)(int);

Write_callback real_write = NULL;
Read_callback real_read = NULL;
Connect_callback real_connect = NULL;
Close_callback real_close = NULL;

void callbacks_init() {
    static bool called = false;
    if (!called) {
        called = true;
    } else {
        return;
    }

    void *h_dl = RTLD_NEXT;
    if (h_dl == NULL) {
        exit(66);
    }

    real_write = (Write_callback)dlsym(h_dl, "write");
    real_read = (Read_callback)dlsym(h_dl, "read");
    real_connect = (Connect_callback)dlsym(h_dl, "connect");
    real_close = (Close_callback)dlsym(h_dl, "close");
//    dlclose(h_dl);
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

bool is_valid(const bson_t* bson);

int establish_ipc() {
    int ipc_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (ipc_sock_fd == -1) {
        perror("Failed to open tmp socket");
        return -1;
    }

    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, IPC_SOCK_PATH);

    int tmp_sock_connect_res = real_connect(ipc_sock_fd, (const struct sockaddr*)&name, sizeof(name));
    if (tmp_sock_connect_res == -1) {
        perror("Connect() tmp socket failed");
        real_close(ipc_sock_fd);
        return -1;
    }

    return ipc_sock_fd;
}

SO_EXPORT int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen){
    callbacks_init();

    if (!is_internet_socket(sock_fd)) {
        return real_connect(sock_fd, addr, addrlen);
    }

    struct sockaddr_in* sin = (struct sockaddr_in*)addr;

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "connect");
    BSON_APPEND_INT32(&bson_request, "sock_fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "port", ntohs(sin->sin_port));
    BSON_APPEND_INT32(&bson_request, "ip", sin->sin_addr.s_addr);

    int ipc_sock_fd = establish_ipc();
    if (ipc_sock_fd == -1) {
        real_write(2, "Failed to establish IPC\n", 24);
        return -1;
    }
    
    int bytes_written = real_write(ipc_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (bytes_written == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

//    real_write(2, "connect\n", 8);

    bson_reader_t* reader = bson_reader_new_from_fd(ipc_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
//    bson_reader_destroy(reader);

    real_write(2, bson_get_data(bson_response), bson_request.len);

    if(!is_valid(bson_response)){
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
        perror("Failed to parse bson: 'result code' is not int32");
        return -1;
    }

    res = bson_iter_int32(&result_code);
    if (-1 == real_close(ipc_sock_fd)){
        perror("Close() tmp socket failed");
    }

//    bson_destroy(&bson_response);
    bson_reader_destroy(reader);
    return res;
}

SO_EXPORT ssize_t read(int sock_fd, void *buf, size_t count){
    callbacks_init();

    if (!is_internet_socket(sock_fd)){
        return real_read(sock_fd, buf, count);
    }

    int ipc_sock_fd = establish_ipc();
    if (ipc_sock_fd == -1) {
        real_write(2, "Failed to establish IPC\n", 24);
        return -1;
    }

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "read");
    BSON_APPEND_INT32(&bson_request, "fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "bytes_to_read", count);

    int bytes_written = real_write(ipc_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (bytes_written == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(ipc_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    if (!is_valid(bson_response)){
        return -1;
    }

    bson_iter_t iter;
    bson_iter_t bytes_read;
    bson_iter_t buffer;

    if (!bson_iter_init(&iter, bson_response)){
        perror("Failed to parse bson: bson_iter_init");
        return -1;
    }
    if (!bson_iter_find_descendant(&iter, "buffer", &buffer)){
        perror("Failed to parse bson: can't find 'buffer'");
        return -1;
    }
    if (!bson_iter_find_descendant(&iter, "bytes_read", &bytes_read)){
        perror("Failed to parse bson: can't find 'bytes_read'");
        return -1;
    }
    int n = bson_iter_int32(&bytes_read);
    bson_iter_binary(&buffer, BSON_SUBTYPE_BINARY, &n, (const uint8_t**)&buf);

//    bson_destroy(bson_response);
    bson_reader_destroy(reader);

    if(-1 == real_close(ipc_sock_fd)){
        perror("Close() tmp socket failed");
        return -1;
    }

    return n;
}

SO_EXPORT ssize_t write(int sock_fd, const void *buf, size_t count){
    callbacks_init();

    if (!is_internet_socket(sock_fd)){
        return real_write(sock_fd, buf, count);
    }

    int ipc_sock_fd = establish_ipc();
    if (ipc_sock_fd == -1) {
        real_write(2, "Failed to establish IPC\n", 24);
        return -1;
    }

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "write");
    BSON_APPEND_INT32(&bson_request, "fd", sock_fd);
    BSON_APPEND_BINARY(&bson_request, "buffer", BSON_SUBTYPE_BINARY, buf, count);
    BSON_APPEND_INT32(&bson_request, "bytes_to_write", count);

    int write_res = real_write(ipc_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (write_res == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(ipc_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    if(!is_valid(bson_response)){
        return -1;
    }
    bson_iter_t iter;
    bson_iter_t bytes_written;
    if(!bson_iter_init(&iter, bson_response)){
        perror("Failed to parse bson: bson_iter_init");
        return -1;
    }
    if(!bson_iter_find_descendant(&iter, "bytes_written", &bytes_written)){
        perror("Failed to parse bson: can't find 'bytes_written'");
        return -1;
    }
    ssize_t res = bson_iter_int32(&bytes_written);

//    bson_destroy(bson_response);
    bson_reader_destroy(reader);

    if(-1 == real_close(ipc_sock_fd)){
        perror("Close() tmp socket failed");
        return -1;
    }

    return res;
}

SO_EXPORT int close(int fd){
    callbacks_init();
    
    if (!is_internet_socket(fd)){
        return real_close(fd);
    }

    int ipc_sock_fd = establish_ipc();
    if (ipc_sock_fd == -1) {
        real_write(2, "Failed to establish IPC\n", 24);
        return -1;
    }

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "close");
    BSON_APPEND_INT32(&bson_request, "fd", fd);

    int write_res = real_write(ipc_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (write_res == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(ipc_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    if(!is_valid(bson_response)){
        return -1;
    }

    bson_iter_t iter;
    bson_iter_t close_res;
    if(!bson_iter_init(&iter, bson_response)){
        perror("Failed to parse bson: bson_iter_init");
        return -1;
    }
    if(!bson_iter_find_descendant(&iter, "close_res", &close_res)){
        perror("Failed to parse bson: can't find 'close_res'");
        return -1;
    }
    int res = bson_iter_int32(&close_res);

    return res;
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
            real_write(2, "Response bson is not valid\n", 27);
            return false;
        } 
    return true;
}
