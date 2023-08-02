GO_CC=go build
C_CC=gcc

ifeq ($(VPNCHAINS_LIB_NAME),)
VPNCHAINS_LIB_NAME := vpnchains_inject.so
endif

ifeq ($(VPNCHAINS_EXECUTABLE_NAME),)
VPNCHAINS_EXECUTABLE_NAME := vpnchains
endif

ifeq ($(VPNCHAINS_OUTPUT_DIR),)
VPNCHAINS_OUTPUT_DIR := build
endif

ifeq ($(LIBBSON_INCLUDE_DIR),)
LIBBSON_INCLUDE_DIR := /usr/include/libbson-1.0
endif


all: pre app lib
pre:
	mkdir -p $(VPNCHAINS_OUTPUT_DIR)
app: pre
	$(GO_CC) -o $(VPNCHAINS_OUTPUT_DIR)/$(VPNCHAINS_EXECUTABLE_NAME) gopkg/main/*.go
lib: pre
	$(C_CC) -shared -fPIC -fvisibility=hidden -I$(LIBBSON_INCLUDE_DIR) -o $(VPNCHAINS_OUTPUT_DIR)/$(VPNCHAINS_LIB_NAME) cpkg/lib.c -lbson-1.0
test: pre
	$(C_CC) -o $(VPNCHAINS_OUTPUT_DIR)/test test/example.c
