package utils

import (
	"fmt"
	"os"
)

func Eprintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}

func Eprintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}
