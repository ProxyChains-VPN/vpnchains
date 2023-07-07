package main

import (
	"abobus/gopkg/ipc"
	"log"
	"net"
)

const DefaultSockAddr = "/tmp/vpnchains.socket"

func errorMsg(path string) string {
	return "Usage: " + path +
		" <lib> <command> [command args...]"
}

func main() {
	conn := ipc.NewConnection(DefaultSockAddr)
	fun := func(conn net.Conn) {
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
	err := conn.Listen(fun)
	if err != nil {
		log.Fatalln(err)
	}

	//args := os.Args
	//if len(args) < 3 {
	//	fmt.Println(errorMsg(args[0]))
	//	os.Exit(0)
	//}

	//libPath := args[1]
	//commandPath := args[2]
	//commandArgs := args[3:]

	//cmd := ipc.CreateCommandWithInjectedLibrary(libPath, commandPath, commandArgs)

	//err := cmd.Run()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//err = cmd.Wait()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//doc, _ := bson.Marshal(
	//	bson.D{
	//		{"call", "write"},
	//		{"Fd", 6},
	//		{"Buffer", []byte("anime")},
	//		{"BytesToWrite", 10050},
	//	},
	//)
	//
	//val, err := ipc.HandleRequest(doc)
	//if err != nil {
	//	log.Fatalln(err)
	//} else {
	//	var res ipc.WriteRequest
	//	bson.Unmarshal(val, &res)
	//	log.Fatalln(res)
	//}

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//go func() {
	//	<-c
	//	os.Remove(DefaultSockAddr)
	//	os.Exit(1)
	//}()
}
