package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	fmt.Println(Reverse("Hello, OTUS!"))
}

func Reverse(str string) string {
	return stringutil.Reverse(str)
}
