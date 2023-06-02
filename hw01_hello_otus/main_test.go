package main

import (
	"testing"
)

func TestReverse(t *testing.T) {
	str := "Hello"
	needStr := "olleH"

	result := Reverse(str)
	if needStr != result {
		t.Error("Работаем не правильно!", result)
	}
}
