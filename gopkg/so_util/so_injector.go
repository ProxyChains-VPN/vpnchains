package so_util

import (
	"os"
	"os/exec"
)

func CreateCommandWithInjectedLibrary(libPath string, commandPath string, args []string) *exec.Cmd {
	cmd := exec.Command(commandPath, args...)

	cmd.Env = append(cmd.Environ(), "LD_PRELOAD="+libPath)
	cmd.Dir = ""

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
