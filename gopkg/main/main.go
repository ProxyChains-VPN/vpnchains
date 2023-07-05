package main

import (
	"abobus/gopkg/so_util"
	"fmt"
	"os"
)

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

	libPath := args[1]
	commandPath := args[2]
	commandArgs := args[3:]

	cmd := so_util.CreateCommandWithInjectedLibrary(libPath, commandPath, commandArgs)
	cmd.Wait()
}
