pragma once
#include <inttypes.h>
#include <sys/socket.h>

typedef size_t (*Read_callback)(int, void*, size_t);
typedef int (*Write_callback)(int, void*, size_t);
typedef int (*Close_callback)(int);

char* read_name = "read";
char* write_name = "write";
char* close_name = "close";

void* get_hDl();

int connect(int, const struct sockaddr*, socklen_t);
size_t read(int, void*, size_t);
int write(int, void*, size_t);