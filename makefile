GO_CC=go build
C_CC=gcc
OUTPUT_DIR=build


all: main lib

main:
	$(GO_CC) -o $(OUTPUT_DIR)/app gopkg/main/main.go
lib :
	$(C_CC) -O2 -shared -fPIC -fvisibility=hidden -o $(OUTPUT_DIR)/string.so cpkg/lib.c
