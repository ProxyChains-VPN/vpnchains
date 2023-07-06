#pragma once
#include <inttypes.h>
#include <sys/socket.h>

typedef int32_t (*Func)();

char* read_name = "read";
char* write_name = "write";
char* close_name = "close";

Func get_real_func(char*);

int connect(int, const struct sockaddr*, socklen_t);
size_t read(int, void*, size_t);
int write(int, void*, size_t);