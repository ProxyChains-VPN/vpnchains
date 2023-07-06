#include <fcntl.h>
#include <dlfcn.h>
#include "lib.h"

void* get_hDl(){
    return dlopen("unistd.h", RTLD_LAZY);
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

    void *hDl = get_hDl();

    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &sock_fd, sizeof(int));
    real_write(tmp_sock_fd, (void*)addr, addrlen);
    real_close(tmp_sock_fd);

    return 0;
}

size_t read(int sock_fd, void *buf, size_t count){
    int tmp_sock_fd;
    int n;

    Read_callback real_read = get_real_read();
    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &sock_fd, sizeof(int));
    n = real_read(tmp_sock_fd, buf, count);
    real_close(tmp_sock_fd);

    return n;
}

int write(int sock_fd, void *buf, size_t count){
    int tmp_sock_fd;

    Write_callback real_write = get_real_write();
    Close_callback real_close = get_real_close();

    tmp_sock_fd = open("/tmp/vpnchains.socket", O_RDWR);
    real_write(tmp_sock_fd, &sock_fd, sizeof(int));
    real_write(tmp_sock_fd, buf, count);
    real_close(tmp_sock_fd);

    return 0;
}
