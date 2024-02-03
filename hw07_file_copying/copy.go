package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFileNoFound           = errors.New("file not found")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileIn, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла: %w", err)
	}

	filiInInfo, err := fileIn.Stat()
	if err != nil {
		return fmt.Errorf("ошибка при получении данных о файле: %w", err)
	}

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("такого файла нет: %w", ErrFileNoFound)
		}
		return fmt.Errorf("ошибка при получении файла: %w", ErrUnsupportedFile)
	}

	fileSize := filiInInfo.Size()

	if fileSize == 0 {
		return fmt.Errorf("такой файл не поддерживается: %w", ErrUnsupportedFile)
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 && offset == 0 {
		fileOut, _ := os.Create(toPath)

		bar := pb.StartNew(int(fileSize))
		bar.SetRefreshRate(time.Millisecond * 100)
		bar.Set(pb.Bytes, true)
		barReader := bar.NewProxyReader(fileIn)

		_, err := io.Copy(fileOut, barReader)
		if err != nil {
			return fmt.Errorf("ошибка при копировании: %w", err)
		}
		bar.Finish()
		return nil
	}

	limitSize := fileSize
	if limit > 0 && limit < fileSize {
		limitSize = limit
	}

	bufferSize := limitSize

	if offset > 0 {
		bufferSize += offset
	}

	buffer := make([]byte, bufferSize)

	barRead := pb.StartNew(int(bufferSize))
	barRead.SetRefreshRate(time.Millisecond * 100)
	barRead.Set(pb.Bytes, true)

	barReader := barRead.NewProxyReader(fileIn)

	var tick int64
	for tick < bufferSize {
		read, err := barReader.Read(buffer[tick:])
		tick += int64(read)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("ошибка при чтении: %w", err)
		}
	}

	errWrite := os.WriteFile(toPath, buffer[offset:tick], os.FileMode(0o644))

	if errWrite != nil {
		return fmt.Errorf("ошибка при записи: %w", errWrite)
	}

	barRead.Finish()

	if err := fileIn.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии файла: %w", err)
	}

	return nil
}
