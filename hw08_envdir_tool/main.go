package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("not found env folder")
		return
	}

	envs, err := ReadDir(args[1])
	if err != nil {
		fmt.Printf("ReadDir: %s", err)
		return
	}

	RunCmd(args, envs)
}
