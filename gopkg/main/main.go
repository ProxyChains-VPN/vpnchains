package main

import (
	"abobus/gopkg/ipc"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const DefaultSockAddr = "/tmp/vpnchains.socket"
const InjectedLibPath = "/usr/lib/libvpnchains_inject.so"

func errorMsg(path string) string {
	return "Usage: " + path +
		" <command> [command args...]"
}

func handleIpc() {
	_ = os.Remove(DefaultSockAddr)

	conn := ipc.NewConnection(DefaultSockAddr)
	handler := func(conn net.Conn) {
		var requestBuf []byte
		_, err := conn.Read(requestBuf)
		if err != nil {
			log.Fatalln(err)
		}

		responseBuf, err := ipc.HandleRequest(requestBuf)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = conn.Write(responseBuf)
		if err != nil {
			log.Fatalln(err)
		}
	}
	err := conn.Listen(handler)
	if err != nil {
		log.Fatalln(err)
	}
}

func sigintHandlerGoroutine() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(DefaultSockAddr)
		os.Exit(1)
	}()
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(errorMsg(args[0]))
		os.Exit(0)
	}

	commandPath := args[1]
	commandArgs := args[2:]

	cmd := ipc.CreateCommandWithInjectedLibrary(InjectedLibPath, commandPath, commandArgs)

	go handleIpc()
	sigintHandlerGoroutine()

	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
