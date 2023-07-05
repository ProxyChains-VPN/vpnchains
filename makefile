CC=go build
#LIB_BUILDMODE=-buildmode=c-shared

all: main

main:
	$(CC) -o build/app pkg/main/main.go