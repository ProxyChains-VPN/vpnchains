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


all: app lib

app:
	$(GO_CC) -o $(VPNCHAINS_OUTPUT_DIR)/$(VPNCHAINS_EXECUTABLE_NAME) gopkg/main/main.go
lib:
	$(C_CC) -O2 -shared -fPIC -fvisibility=hidden -o $(VPNCHAINS_OUTPUT_DIR)/$(VPNCHAINS_LIB_NAME) cpkg/lib.c
