package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Get envs OK", func(t *testing.T) {
		cmdString := []string{
			"./go-envdir",
			"testdata/env",
			"/bin/bash",
			"testdata/echo.sh",
			"arg1=1",
			"arg2=2",
		}

		envs := Environment(map[string]EnvValue{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: `"hello"`, NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		})

		code := RunCmd(cmdString, envs)

		require.Equal(t, 0, code)
	})
	t.Run("Should be more arguments Fail", func(t *testing.T) {
		cmdString := []string{
			"./go-envdir",
			"testdata/env",
		}

		envs := Environment(map[string]EnvValue{
			"BAR": {Value: "bar", NeedRemove: false},
		})

		code := RunCmd(cmdString, envs)

		require.Equal(t, 1, code)
	})
	t.Run("Command exit 2 Fail", func(t *testing.T) {
		cmdString := []string{
			"./go-envdir",
			"testdata/env",
			"/bin/bash",
			"testdata/error.sh",
		}

		envs := Environment(map[string]EnvValue{
			"BAR": {Value: "bar", NeedRemove: false},
		})

		code := RunCmd(cmdString, envs)

		require.Equal(t, 2, code)
	})
}
