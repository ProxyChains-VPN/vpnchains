#pragma once

#define SO_EXPORT __attribute__((visibility("default")))
#define GET_HDL() dlopen(LIBC_SO, RTLD_LAZY)