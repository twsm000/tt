package utils

import (
	"fmt"
	"os"
)

func MustGet[T any](t T, err error) T {
	TryTerminate(err)
	return t
}

func TryTerminate(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
