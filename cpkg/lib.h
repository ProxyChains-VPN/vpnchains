#pragma once
#include <inttypes.h>
#include <sys/socket.h>
#include <bson/bson.h>

typedef size_t (*Read_callback)(int, void*, size_t);
typedef int (*Write_callback)(int, void*, size_t);
typedef int (*Close_callback)(int);

void* get_hDl();

int connect(int, const struct sockaddr*, socklen_t);
ssize_t read(int, void*, size_t);
ssize_t write(int, const void*, size_t);