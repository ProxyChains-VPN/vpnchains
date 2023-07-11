package main

import (
	"abobus/gopkg/ipc"
	"abobus/gopkg/vpn"
	"fmt"
	"log"
	"net"
	"os"
)

const DefaultSockAddr = "/tmp/vpnchains.socket"
const InjectedLibPath = "/usr/lib/libvpnchains_inject.so"
const BufSize = 100500

func errorMsg(path string) string {
	return "Usage: " + path +
		" <command> [command args...]"
}

func handleIpc(ready chan struct{}) {
	err := os.Remove(DefaultSockAddr)
	if err != nil {
		log.Fatalln(err)
	}

	var buf = make([]byte, BufSize)

	conn := ipc.NewConnection(DefaultSockAddr)
	handler := func(conn net.Conn) {
		n, err := conn.Read(buf)
		requestBuf := buf[:n]

		if err != nil {
			log.Fatalln(err)
		}

		responseBuf, err := ipc.HandleRequest(requestBuf)
		if responseBuf == nil && err != nil {
			log.Fatalln(err) // вроде как невозможно
		} else if err != nil {
			log.Println(err, ". Returning error response.")
		}

		_, err = conn.Write(responseBuf)
		if err != nil {
			log.Println(err)
		}
	}

	ready <- struct{}{}
	err = conn.Listen(handler)
	if err != nil {
		log.Println("sldfadsf")
		log.Fatalln(err)
	}
}

//func sigintHandlerGoroutine() {
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//	go func() {
//		<-c
//		os.Remove(DefaultSockAddr)
//		os.Exit(1)
//	}()
//}

func main() {
	//args := os.Args
	//if len(args) < 2 {
	//	fmt.Println(errorMsg(args[0]))
	//	os.Exit(0)
	//}
	//
	//commandPath := args[1]
	//commandArgs := args[2:]
	//
	//cmd := ipc.CreateCommandWithInjectedLibrary(InjectedLibPath, commandPath, commandArgs)
	//
	//ready := make(chan struct{})
	//go handleIpc(ready)
	//
	//<-ready
	//err := cmd.Start()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//err = cmd.Wait()
	//if err != nil {
	//	log.Fatalln(83, err)
	//}

	res, err := vpn.WireguardConfigFromFile("wg0.conf")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.ToString())
}
