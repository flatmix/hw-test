package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	inFilePath := "testdata/input.txt"
	invalidInFilePath := "testdata/input1.txt"
	outFilePath := "testdata/output.txt"

	t.Run("full copy file", func(t *testing.T) {
		err := Copy(inFilePath, outFilePath, 0, 0)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
		fileOut, err := os.Open(outFilePath)

		require.NoError(t, err)

		filiOutInfo, err := fileOut.Stat()
		if err != nil {
			t.Fatalf("Ошибка при получении данных о файле: %s", err)
		}
		fileSize := filiOutInfo.Size()

		require.Equal(t, int64(6617), fileSize)
	})
	t.Run("fail offset", func(t *testing.T) {
		err := Copy(inFilePath, outFilePath, 7000, 0)
		require.Error(t, ErrOffsetExceedsFileSize, err)
	})
	t.Run("file not found", func(t *testing.T) {
		err := Copy(invalidInFilePath, outFilePath, 0, 0)
		require.Error(t, ErrFileNoFound, err)
	})
	t.Run("file without", func(t *testing.T) {
		err := Copy("/dev/urandom", outFilePath, 0, 7000)
		require.Error(t, ErrUnsupportedFile, err)
	})
	t.Run("limit more filesize", func(t *testing.T) {
		err := Copy(inFilePath, outFilePath, 0, 7000)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
		fileOut, err := os.Open(outFilePath)

		require.NoError(t, err)

		filiOutInfo, err := fileOut.Stat()
		if err != nil {
			t.Fatalf("Ошибка при получении данных о файле: %s", err)
		}
		fileSize := filiOutInfo.Size()

		require.Equal(t, int64(6617), fileSize)
	})

	err := os.Remove(outFilePath)
	require.NoError(t, err)
}
