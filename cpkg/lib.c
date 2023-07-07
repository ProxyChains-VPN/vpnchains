#include <fcntl.h>
#include <dlfcn.h>
#include <netinet/in.h>
#include <sys/stat.h>
#include "lib.h"

void* get_hDl(){
    char* lib_name = "unistd.h";
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

Close_callback get_real_close(){
    void *hDl = get_hDl();
    Close_callback real_close = (Close_callback)dlsym(hDl, "close");
    dlclose(hDl);
    return real_close;
}

int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen){

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);  
    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");
    connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));

    void *hDl = get_hDl();
    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    struct sockaddr_in* sin = (struct sockaddr_in*)addr;
    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "Call", "connect");
    BSON_APPEND_INT32(&bson_request, "SockFd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "Port", ntohs(sin->sin_port));
    BSON_APPEND_INT32(&bson_request, "Ip", sin->sin_addr.s_addr);
    
    real_write(tmp_sock_fd, &bson_request, bson_request.len);

    bson_reader_t* reader = bson_reader_new_from_fd(tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    bson_iter_t iter;
    bson_iter_t result_code;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "ResultCode", &result_code);
    int res = bson_iter_int32(&result_code);

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

    int tmp_sock_fd = socket(AF_UNIX, SOCK_STREAM, 0);  
    struct sockaddr_un name;
    memset(&name, 0, sizeof(name));
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, "/tmp/vpnchains.socket");
    connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "Call", "read");
    BSON_APPEND_INT32(&bson_request, "Fd", sock_fd);
    BSON_APPEND_INT32(&bson_request, "BytesToRead", count);

    real_write(tmp_sock_fd, &bson_request, bson_request.len);

    bson_reader_t* reader = bson_reader_new_from_fd (tmp_sock_fd, false);
    const bson_t* bson_response = bson_reader_read(reader, NULL);
    bson_iter_t bytes_read;
    bson_iter_init(&iter, bson_response);
    bson_iter_find_descendant(&iter, "BytesRead", &bytes_read);
    int n = bson_iter_int32(&bytes_read);
    /*read*/

    real_close(tmp_sock_fd);

    return n;
}

ssize_t write(int sock_fd, const void *buf, size_t count){

    Write_callback real_write = get_real_write();

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
    connect(tmp_sock_fd, (const struct sockaddr*)&name, sizeof(name));

    bson_t bson_request = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&bson_request, "Call", "write");
    BSON_APPEND_INT32(&bson_request, "Fd", sock_fd);
    BSON_APPEND_BINARY(&bson_request, "Buffer", BSON_SUBTYPE_BINARY, buf, count);
    BSON_APPEND_INT32(&bson_request, "BytesToWrite", count);

    real_write(tmp_sock_fd, &bson_request, bson_request.len);

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