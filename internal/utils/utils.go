package utils

import (
	"fmt"
	"io"
)

func PrintFromReader(f io.Reader, bufferSize uint) {
	for {
		bytes := make([]byte, 64)
		_, err := f.Read(bytes)
		if err == io.EOF {
			return
		}
		fmt.Printf("%s", bytes)
	}
}
