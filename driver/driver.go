package driver

import (
	"fmt"
	"os"
)

type Options struct {
	Out string
}

func Compile(path string, opt Options) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Printf("File contents:\n%s", src);

	return nil
}
