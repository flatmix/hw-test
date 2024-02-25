package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func removeEnv(name string) error {
	_, ok := os.LookupEnv(name)
	if ok {
		err := os.Unsetenv(name)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, envs Environment) (returnCode int) {
	var commandExec string

	if len(cmd) < 3 {
		fmt.Printf("should be more arguments in the cmd, now: %d \n", len(cmd))
		return 1
	}

	commandLine := make([]string, 0, len(cmd)-3)
	var err error

	for name, env := range envs {
		err = removeEnv(name)
		if err != nil {
			fmt.Printf("removeEnv for %s, error: %v \n", name, err)
			return 1
		}
		if !env.NeedRemove {
			err = os.Setenv(name, env.Value)
			if err != nil {
				fmt.Printf("setenv for %s, value: %s, error: %v \n", name, env.Value, err)
				return 1
			}
		}
	}

	commandExec = cmd[2]
	if len(cmd) > 3 {
		commandLine = append(commandLine, cmd[3:]...)
	}

	command := exec.Command(commandExec, commandLine...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}

	return
}
