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
#include <libbson-1.0/bson/bson.h>

typedef size_t (*Read_callback)(int, void*, size_t);
typedef int (*Write_callback)(int, const void*, size_t);
typedef int (*Connect_callback)(int, const struct sockaddr*, socklen_t);
typedef int (*Close_callback)(int);

bool is_internet_socket(int sock_fd){
    struct stat statbuf;
    fstat(fd, &statbuf);
    if (!S_ISSOCK(statbuf.st_mode)){
        return false;
    }

    struct sockaddr addr;
    getsockname(fd, &addr, sizeof(addr));
    return addr.sa_family != AF_UNIX;
}

int is_valid(const bson_t* bson);

void establish_ipc() {
    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (tmp_sock_fd == -1) {
        perror("Failed to open tmp socket");
        return -1;
    }

    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");

    int tmp_sock_connect_res = real_connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));
    if (tmp_sock_connect_res == -1) {
        perror("Connect() tmp socket failed");
        return -1;
    }
}

Write_callback get_real_write(){
    void *hDl = GET_HDL();
    Write_callback real_write = (Write_callback)dlsym(hDl, "write");
    dlclose(hDl);
    return real_write;
}

Read_callback get_real_read(){
    void *hDl = GET_HDL();
    Read_callback real_read = (Read_callback)dlsym(hDl, "read");
    dlclose(hDl);
    return real_read;
}

Connect_callback get_real_connect(){
    void *hDl = GET_HDL();
    Connect_callback real_connect = (Connect_callback)dlsym(hDl, "connect");
    dlclose(hDl);
    return real_connect;
}

Close_callback get_real_close(){
    void *hDl = GET_HDL();
    Close_callback real_close = (Close_callback)dlsym(hDl, "close");
    dlclose(hDl);
    return real_close;
}

SO_EXPORT int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen){
    Connect_callback real_connect = get_real_connect();
    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    if (!is_internet_socket(sock_fd)){
        return real_connect(sock_fd, addr, addrlen);
    }

    struct sockaddr_in* sin = (struct sockaddr_in*)addr;
    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "connect");
    BSON_APPEND_INT32(&bson_request, "sock_fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "port", ntohs(sin->sin_port));
    BSON_APPEND_INT32(&bson_request, "ip", sin->sin_addr.s_addr);
    
    int bytes_written = real_write(tmp_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (bytes_written == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
//    bson_reader_destroy(reader);

    real_write(2, bson_get_data(bson_response), bson_request.len);

    if(!is_valid(bson_response)){
        return -1;
    }

    int res = -1;

    bson_iter_t iter;
    bson_iter_t result_code;
    if(!bson_iter_init(&iter, bson_response)){
        perror("Failed to parse bson: bson_iter_init");
        return -1;
    }

    if(!bson_iter_find_descendant(&iter, "result_code", &result_code)){
        perror("Failed to parse bson: can't find 'result_code'");
        return -1;
    }

    if(!BSON_ITER_HOLDS_INT32(&result_code)){
        perror("Failed to parse bson: 'result code' is not int32");
        return -1;
    }

    res = bson_iter_int32(&result_code);
    // real_write(2, "and we a re here\n", 17);
    if(-1 == real_close(tmp_sock_fd)){
        perror("Close() tmp socket failed");
    }

//    bson_destroy(&bson_response);
    bson_reader_destroy(reader);
    return res;
}

SO_EXPORT ssize_t read(int sock_fd, void *buf, size_t count){
    Read_callback real_read = get_real_read();

    if (!is_internet_socket(sock_fd)){
        return real_read(sock_fd, buf, count);
    }

    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

//    real_write(2, "abobaREAD\n", 11);

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (tmp_sock_fd == -1) {
        perror("Failed to open tmp socket");
        return -1;
    }

    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");

    Connect_callback real_connect = get_real_connect();
    int tmp_sock_connect_res = real_connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));
    if (tmp_sock_connect_res == -1) {
        perror("Connect() tmp socket failed");
        return -1;
    }

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "read");
    BSON_APPEND_INT32(&bson_request, "fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "bytes_to_read", count);

    int bytes_written = real_write(tmp_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (bytes_written == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    if(!is_valid(bson_response)){
        return -1;
    }
    bson_iter_t iter;
    bson_iter_t bytes_read;
    bson_iter_t buffer;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "bytes_read", &bytes_read);
    bson_iter_find_descendant(&iter, "buffer", &buffer);
    int n = bson_iter_int32(&bytes_read);
    bson_iter_binary(&buffer, BSON_SUBTYPE_BINARY, &n, (const uint8_t**)&buf);

//    bson_destroy(bson_response);
    bson_reader_destroy(reader);

    if(-1 == real_close(tmp_sock_fd)){
        perror("Close() tmp socket failed");
        return -1;
    }

    return n;
}

SO_EXPORT ssize_t write(int sock_fd, const void *buf, size_t count){
    Write_callback real_write = get_real_write();

    if (!is_internet_socket(sock_fd)){
        return real_write(sock_fd, buf, count);
    }

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (tmp_sock_fd == -1) {
        perror("Failed to open tmp socket");
        return -1;
    }

    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");

    Connect_callback real_connect = get_real_connect();
    int tmp_sock_connect_res = real_connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));
    if (tmp_sock_connect_res == -1) {
        perror("Connect() tmp socket failed");
        return -1;
    }

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "write");
    BSON_APPEND_INT32(&bson_request, "fd", sock_fd);
    BSON_APPEND_BINARY(&bson_request, "buffer", BSON_SUBTYPE_BINARY, buf, count);
    BSON_APPEND_INT32(&bson_request, "bytes_to_write", count);

    int write_res = real_write(tmp_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (write_res == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    if(!is_valid(bson_response)){
        return -1;
    }
    bson_iter_t iter;
    bson_iter_t bytes_written;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "bytes_written", &bytes_written);
    ssize_t res = bson_iter_int32(&bytes_written);

//    bson_destroy(bson_response);
    bson_reader_destroy(reader);

    Close_callback real_close = get_real_close();
    if(-1 == real_close(tmp_sock_fd)){
        perror("Close() tmp socket failed");
        return -1;
    }

    return res;
}

SO_EXPORT int close(int fd){
    Close_callback real_close = get_real_close();

    struct stat statbuf;
    fstat(fd, &statbuf);
    if(!S_ISSOCK(statbuf.st_mode)){
        return real_close(fd);
    }

    struct sockaddr addr;
    getsockname(fd, &addr, sizeof(addr));
    if (addr.sa_family == AF_UNIX) {
        return real_close(fd);
    }

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (tmp_sock_fd == -1) {
        perror("Failed to open tmp socket");
        return -1;
    }

    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");

    Connect_callback real_connect = get_real_connect();
    int tmp_sock_connect_res = real_connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));
    if (tmp_sock_connect_res == -1) {
        perror("Connect() tmp socket failed");
        return -1;
    }

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "call", "close");
    BSON_APPEND_INT32(&bson_request, "fd", fd);

    int write_res = real_write(tmp_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (write_res == -1) {
        perror("Write() to tmp socket failed");
        return -1;
    }

    bson_destroy(&bson_request);

    bson_reader_t* reader = bson_reader_new_from_fd(tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    if(!is_valid(bson_response)){
        return -1;
    }

    bson_iter_t iter;
    bson_iter_t close_res;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "close_res", &close_res);
    int res = bson_iter_int32(&close_res);

    return res;
}

int is_valid(const bson_t* bson){
    if (!bson_validate(
        bson, 
        BSON_VALIDATE_UTF8 
        | BSON_VALIDATE_DOLLAR_KEYS 
        | BSON_VALIDATE_DOT_KEYS 
        | BSON_VALIDATE_UTF8_ALLOW_NULL 
        | BSON_VALIDATE_EMPTY_KEYS,
        NULL)) {
            real_write(2, "Response bson is not valid\n", 27);
            return -1;
        } 
    return 0;
}
