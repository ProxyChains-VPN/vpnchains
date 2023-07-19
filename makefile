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


all: pre app lib
pre:
	mkdir -p $(VPNCHAINS_OUTPUT_DIR)
app: pre
	$(GO_CC) -o $(VPNCHAINS_OUTPUT_DIR)/$(VPNCHAINS_EXECUTABLE_NAME) gopkg/main/main.go gopkg/main/ipc.go
lib: pre
	$(C_CC) -shared -fPIC -fvisibility=hidden -o $(VPNCHAINS_OUTPUT_DIR)/$(VPNCHAINS_LIB_NAME) cpkg/lib.c -lbson-1.0
test: pre
	$(C_CC) -o $(VPNCHAINS_OUTPUT_DIR)/test test/example.c
