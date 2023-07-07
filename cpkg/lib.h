#pragma once
#include <inttypes.h>
#include <libbson-1.0/bson/bson.h>

typedef size_t (*Read_callback)(int, void*, size_t);
typedef int (*Write_callback)(int, const void*, size_t);
typedef int (*Connect_callback)(int, const struct sockaddr*, socklen_t);
typedef int (*Close_callback)(int);

void* get_hDl();

int connect(int, const struct sockaddr*, socklen_t);
ssize_t read(int, void*, size_t);
ssize_t write(int, const void*, size_t);

#define SO_VISIBLE __attribute__((visibility("default")))