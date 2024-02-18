package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func createEnvFile(filename, data string) error {
	buffer := []byte(data)
	errWrite := os.WriteFile(filename, buffer, os.FileMode(0o644))
	if errWrite != nil {
		return fmt.Errorf("create and write file: %w", errWrite)
	}

	return nil
}

func removeEnvFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("remove file: %w", err)
	}

	return nil
}

func TestReadDir(t *testing.T) {
	envPath := "testdata/env"

	t.Run("Get envs OK", func(t *testing.T) {
		envs, err := ReadDir(envPath)

		require.NoError(t, err)

		envsExp := Environment(map[string]EnvValue{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: `"hello"`, NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		})

		require.Equal(t, envsExp, envs)
	})
	t.Run("Add file envs OK", func(t *testing.T) {
		err := createEnvFile("testdata/env/FOR_UNIT_TEST", "test data")
		require.NoError(t, err)

		envs, err := ReadDir(envPath)
		require.NoError(t, err)

		envsExp := Environment(map[string]EnvValue{
			"BAR":           {Value: "bar", NeedRemove: false},
			"EMPTY":         {Value: "", NeedRemove: false},
			"FOO":           {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO":         {Value: `"hello"`, NeedRemove: false},
			"UNSET":         {Value: "", NeedRemove: true},
			"FOR_UNIT_TEST": {Value: "test data", NeedRemove: false},
		})

		require.Equal(t, envsExp, envs)

		err = removeEnvFile("testdata/env/FOR_UNIT_TEST")
		require.NoError(t, err)
	})
	t.Run("Add file envs with incorrect filename ", func(t *testing.T) {
		err := createEnvFile("testdata/env/FOR_UNIT=_TEST_SECOND", "test data second")
		require.NoError(t, err)

		_, err = ReadDir(envPath)
		require.ErrorIs(t, err, ErrorFilename)

		err = removeEnvFile("testdata/env/FOR_UNIT=_TEST_SECOND")
		require.NoError(t, err)
	})
	t.Run("Incorrect path Fail", func(t *testing.T) {
		incorrectEnvPath := "testdata_sec/env"

		_, err := ReadDir(incorrectEnvPath)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}
