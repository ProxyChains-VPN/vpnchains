#include <fcntl.h>
#include <dlfcn.h>
#include <netinet/in.h>
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
    int tmp_sock_fd;
    struct sockaddr_in* sin = (struct sockaddr_in*)addr;

    void *hDl = get_hDl();

    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    bson_t b = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&b, "Call", "connect");
    BSON_APPEND_INT32(&b, "SockFd", sock_fd);
    BSON_APPEND_INT32(&b, "Port", sin->sin_port);
    BSON_APPEND_INT32(&b, "Ip", sin->sin_addr.s_addr);

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &b, b.len);
    real_close(tmp_sock_fd);

    return 0;
}

ssize_t read(int sock_fd, void *buf, size_t count){
    int tmp_sock_fd;
    int n;

    Read_callback real_read = get_real_read();
    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    bson_t b = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&b, "Call", "read");
    BSON_APPEND_INT32(&b, "Fd", sock_fd);
    BSON_APPEND_INT32(&b, "BytesToRead", count);

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &b, b.len);
    //n = real_read(tmp_sock_fd, buf, count);
    real_close(tmp_sock_fd);

    return n;
}

ssize_t write(int sock_fd, const void *buf, size_t count){
    int tmp_sock_fd;

    Write_callback real_write = get_real_write();
    Close_callback real_close =get_real_close();

    bson_t b = BSON_INITIALIZER;
    BSON_APPEND_UTF8(&b, "Call", "write");
    BSON_APPEND_INT32(&b, "Fd", sock_fd);
    BSON_APPEND_BINARY(&b, "Buffer", BSON_SUBTYPE_BINARY, buf, count);
    BSON_APPEND_INT32(&b, "BytesToWrite", count);

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &b, b.len);
    real_close(tmp_sock_fd);

    return 0;
}
