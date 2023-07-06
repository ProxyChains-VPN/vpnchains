package main

import (
	"fmt"
	"os"
)

const DefaultSockAddr = "/tmp/vpnchains.socket"

func errorMsg(path string) string {
	return "Usage: " + path +
		" <lib> <command> [command args...]"
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println(errorMsg(args[0]))
		os.Exit(0)
	}

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
}
