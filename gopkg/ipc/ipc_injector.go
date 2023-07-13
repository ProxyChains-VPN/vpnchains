package ipc

import (
	"os"
	"os/exec"
)

// CreateCommandWithInjectedLibrary creates a new exec.Cmd instance with an injected library.
// In other words, it inserts a library whose path is libPath into the command's LD_PRELOAD environment variable,
// and returns the command.
// libPath - path to the library
// commandPath - path to the command
// args - command arguments
func CreateCommandWithInjectedLibrary(libPath string, commandPath string, args []string) *exec.Cmd {
	cmd := exec.Command(commandPath, args...)

	cmd.Env = append(cmd.Environ(), "LD_PRELOAD="+libPath)
	cmd.Dir = ""

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
