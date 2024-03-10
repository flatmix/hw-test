package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrorFilename = errors.New("incorrect env filename")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environments := Environment{}
	envFiles, err := os.ReadDir(dir)
	if err != nil {
		return Environment{}, fmt.Errorf("os readDir: %w", err)
	}
	for _, envFile := range envFiles {
		envFileInfo, fileErr := envFile.Info()

		if fileErr != nil {
			return Environment{}, fmt.Errorf("file.Info: %w", fileErr)
		}
		name := envFile.Name()
		if strings.Contains(name, "=") {
			return Environment{}, ErrorFilename
		}
		if envFileInfo.Size() == 0 {
			environments[name] = EnvValue{Value: "", NeedRemove: true}

			continue
		}

		file, err := os.Open(fmt.Sprintf("%s/%s", dir, name))
		if err != nil {
			return Environment{}, fmt.Errorf("file open: %w", err)
		}

		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)
		var lines [][]byte
		for fileScanner.Scan() {
			lines = append(lines, fileScanner.Bytes())
		}
		err = file.Close()
		if err != nil {
			return Environment{}, fmt.Errorf("file close: %w", err)
		}
		prepareValue := strings.TrimRight(string(bytes.ReplaceAll(lines[0], []byte("\x00"), []byte("\n"))), " \t")
		environments[name] = EnvValue{Value: prepareValue, NeedRemove: false}
	}

	return environments, nil
}
