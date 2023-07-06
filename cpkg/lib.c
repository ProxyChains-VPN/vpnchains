#include <fcntl.h>
#include <dlfcn.h>
#include "lib.h"

Func get_real_func(char* func_name){
    char* lib_name = "unistd.h";
    void *hDl = dlopen(lib_name, RTLD_LAZY);
    Func real_func = (Func)dlsym(hDl, func_name);

    dlclose(hDl);

    return real_func;
}

int connect(int sock_fd, const struct sockaddr *addr, socklen_t addrlen){
    int tmp_sock_fd;

    Func real_write = get_real_func(write_name);
    Func real_close = get_real_func(close_name);

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &sock_fd, sizeof(int));
    real_write(tmp_sock_fd, addr, addrlen);
    real_close(tmp_sock_fd);

    return 0;
}

size_t read(int sock_fd, void *buf, size_t count){
    int tmp_sock_fd;
    int n;

    Func real_write = get_real_func(write_name);
    Func real_read = get_real_func(read_name);
    Func real_close = get_real_func(close_name);

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &sock_fd, sizeof(int));
    n = real_read(tmp_sock_fd, buf, count);
    real_close(tmp_sock_fd);

    return n;
}

int write(int sock_fd , void *buf, size_t count){
    int tmp_sock_fd;

    Func real_write = get_real_func(write_name);
    Func real_close = get_real_func(close_name);

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &sock_fd, sizeof(int));
    real_write(tmp_sock_fd, buf, count);
    real_close(tmp_sock_fd);

    return 0;
}