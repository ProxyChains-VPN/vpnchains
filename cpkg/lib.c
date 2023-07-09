#include <fcntl.h>
#include <dlfcn.h>
#include <netinet/in.h>
#include <sys/stat.h>
#include <sys/un.h>
#include "lib.h"
#include <assert.h>

#include <string.h>
#include <gnu/lib-names.h>

void* get_hDl(){
    char* lib_name = LIBC_SO;
    return dlopen(lib_name, RTLD_LAZY);
}

Write_callback get_real_write(){
    void *hDl = get_hDl();
    Write_callback real_write = (Write_callback)dlsym(hDl, "write");
    dlclose(hDl);
    return real_write;
}

Read_callback get_real_read(){
    void *hDl = get_hDl();
    Read_callback real_read = (Read_callback)dlsym(hDl, "read");
    dlclose(hDl);
    return real_read;
}

Connect_callback get_real_connect(){
    void *hDl = get_hDl();
    Connect_callback real_connect = (Connect_callback)dlsym(hDl, "connect");
    dlclose(hDl);
    return real_connect;
}

Close_callback get_real_close(){
    void *hDl = get_hDl();
    Close_callback real_close = (Close_callback)dlsym(hDl, "close");
    dlclose(hDl);
    return real_close;
}

SO_VISIBLE int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen){
    Connect_callback real_connect = get_real_connect();
    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (tmp_sock_fd == -1) {
        char* err = strerror(errno);
        real_write(2, err, strlen(err));
    }

    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");

    int tmp_sock_connect_res = real_connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));
    if (tmp_sock_connect_res == -1) {
        char* err = strerror(errno);
        real_write(2, err, strlen(err));
    }

    struct sockaddr_in* sin = (struct sockaddr_in*)addr;
    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "Call", "connect");
    BSON_APPEND_INT32(&bson_request, "SockFd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "Port", ntohs(sin->sin_port));
    BSON_APPEND_INT32(&bson_request, "Ip", sin->sin_addr.s_addr);
    
    int bytes_written = real_write(tmp_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (bytes_written == -1) {
        char* err = strerror(errno);
        real_write(2, err, strlen(err));
    }

    bson_reader_t* reader = bson_reader_new_from_fd(tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);

    int res = -1;

    bson_iter_t iter;
    bson_iter_t result_code;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "ResultCode", &result_code);

    if (BSON_ITER_HOLDS_INT32(&result_code)) {
        res = bson_iter_int32(&result_code);
    }
    real_write(2, "and we a re here\n", 17);


//    char buffer[100500];
//    Read_callback real_read = get_real_read();

//    int bytes_read = real_read(tmp_sock_fd, buffer, 100500);
//    if (bytes_read == -1) {
//        char* err = strerror(errno);
//        real_write(2, err, strlen(err));
//    }

//    real_write(2, buffer, bytes_read);


    real_close(tmp_sock_fd);

    return res;
}

ssize_t read(int sock_fd, void *buf, size_t count){

    Read_callback real_read = get_real_read();

    struct stat statbuf;
    fstat(sock_fd, &statbuf);
    if(!S_ISSOCK(statbuf.st_mode)){
        return real_read(sock_fd, buf, count);
    }

    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

//    real_write(2, "abobaREAD\n", 11);

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");

    Connect_callback real_connect = get_real_connect();
    real_connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "Call", "read");
    BSON_APPEND_INT32(&bson_request, "Fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "BytesToRead", count);

    int bytes_written = real_write(tmp_sock_fd, bson_get_data(&bson_request), bson_request.len);
    if (bytes_written == -1) {
        char* err = strerror(errno);
        real_write(2, err, strlen(err));
    } else {
        real_write(2, "ok", 3);
    }

    bson_reader_t* reader = bson_reader_new_from_fd (tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    bson_iter_t iter;
    bson_iter_t bytes_read;
    bson_iter_t buffer;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "BytesRead", &bytes_read);
    bson_iter_find_descendant(&iter, "Buffer", &buffer);
    int n = bson_iter_int32(&bytes_read);
    bson_iter_binary(&buffer, BSON_SUBTYPE_BINARY, &n, (const uint8_t**)&buf);

    real_close(tmp_sock_fd);

    return n;
}

ssize_t write(int sock_fd, const void *buf, size_t count){
    Write_callback real_write = get_real_write();
    real_write(2, "abobaWRIT\n", 11);

    struct stat statbuf;
    fstat(sock_fd, &statbuf);
    if(!S_ISSOCK(statbuf.st_mode)){
        return real_write(sock_fd, buf, count);
    }

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");

    Connect_callback real_connect = get_real_connect();
    real_connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "Call", "write");
    BSON_APPEND_INT32(&bson_request, "Fd", sock_fd);
    BSON_APPEND_BINARY(&bson_request, "Buffer", BSON_SUBTYPE_BINARY, buf, count);
    BSON_APPEND_INT32(&bson_request, "BytesToWrite", count);

    real_write(tmp_sock_fd, bson_get_data(&bson_request), bson_request.len);

    bson_reader_t* reader = bson_reader_new_from_fd(tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    bson_iter_t iter;
    bson_iter_t bytes_written;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "BytesWritten", &bytes_written);
    ssize_t res = bson_iter_int32(&bytes_written);

    Close_callback real_close = get_real_close();
    real_close(tmp_sock_fd);

    return res;
}